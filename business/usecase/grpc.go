// Package usecase provides business logic.
package usecase

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"sync"

	"github.com/pkg/errors"

	"github.com/forest33/warthog/adapter/grpc"
	"github.com/forest33/warthog/business/entity"
	"github.com/forest33/warthog/pkg/logger"
)

// GrpcUseCase object capable of interacting with GrpcUseCase
type GrpcUseCase struct {
	ctx                    context.Context
	log                    *logger.Zerolog
	grpcClient             GrpcClient
	k8sClient              K8SClient
	services               []*entity.Service
	workspaceRepo          WorkspaceRepo
	curServerID            int64
	curConnectedServerID   int64
	curService             string
	curMethod              string
	curServer              *entity.WorkspaceItemServer
	curServerClientOptions []grpc.ClientOpt
	forwardPorts           map[int16]*forwardPort
	muForwardPorts         sync.RWMutex
	infoCh                 chan *entity.Info
	errorCh                chan *entity.Error
}

// GrpcClient is an interface for working with the gRPC client
type GrpcClient interface {
	SetSettings(cfg *entity.Settings)
	Connect(addr string, auth *entity.Auth, opts ...grpc.ClientOpt) error
	AddProtobuf(path ...string)
	AddImport(path ...string)
	LoadFromProtobuf() ([]*entity.Service, []*entity.ProtobufError, *entity.ProtobufError)
	LoadFromReflection() ([]*entity.Service, error)
	GetResponseChannel() chan *entity.QueryResponse
	GetSentCounter() uint
	Query(method *entity.Method, data map[string]interface{}, metadata []string) error
	CancelQuery()
	CloseStream()
	Close()
}

// K8SClient is an interface for working with the K8S
type K8SClient interface {
	PortForward(r *entity.K8SPortForward) (entity.PortForwardControl, error)
}

type forwardPort struct {
	control entity.PortForwardControl
	hash    string
}

// NewGrpcUseCase creates a new GrpcUseCase
func NewGrpcUseCase(ctx context.Context, log *logger.Zerolog, grpcClient GrpcClient, k8sClient K8SClient, workspaceRepo WorkspaceRepo) *GrpcUseCase {
	return &GrpcUseCase{
		ctx:           ctx,
		log:           log,
		grpcClient:    grpcClient,
		k8sClient:     k8sClient,
		workspaceRepo: workspaceRepo,
		infoCh:        make(chan *entity.Info),
		errorCh:       make(chan *entity.Error),
	}
}

// LoadServer reads the server description from the database and returns it to the GUI
func (uc *GrpcUseCase) LoadServer(payload map[string]interface{}) *entity.GUIResponse {
	req := &entity.ServerRequest{}
	if err := req.Model(payload); err != nil {
		return entity.ErrorGUIResponse(err)
	}

	var (
		query  *entity.Workspace
		server *entity.Workspace
		warn   []*entity.ProtobufError
		err    error
	)

	server, err = uc.workspaceRepo.GetByID(req.ID)
	if err != nil {
		uc.log.Error().
			Int64("id", req.ID).
			Msgf("failed to get server: %v", err)
		return entity.ErrorGUIResponse(err)
	}

	if server.Type == entity.WorkspaceTypeQuery {
		query = server
		server, err = uc.workspaceRepo.GetByID(*query.ParentID)
		if err != nil {
			uc.log.Error().
				Int64("id", req.ID).
				Msgf("failed to get server: %v", err)
			return entity.ErrorGUIResponse(err)
		}
	}

	uc.curServer = server.Data.(*entity.WorkspaceItemServer)
	uc.curServerClientOptions = nil
	uc.curServerClientOptions = make([]grpc.ClientOpt, 0, 4)

	if uc.curServer.UseGrpcWeb {
		uc.curServerClientOptions = append(uc.curServerClientOptions, grpc.WithUseWeb())
	}

	if uc.curServer.NoTLS {
		uc.curServerClientOptions = append(uc.curServerClientOptions, grpc.WithNoTLS())
	} else {
		if uc.curServer.Insecure {
			uc.curServerClientOptions = append(uc.curServerClientOptions, grpc.WithInsecure())
		}
		uc.curServerClientOptions = append(uc.curServerClientOptions, grpc.WithRootCertificate(uc.curServer.RootCertificate))
		if uc.curServer.ClientCertificate != "" && uc.curServer.ClientKey != "" {
			uc.curServerClientOptions = append(uc.curServerClientOptions,
				grpc.WithClientCertificate(uc.curServer.ClientCertificate),
				grpc.WithClientKey(uc.curServer.ClientKey))
		}
	}

	if uc.curServer.UseReflection {
		err = uc.connect(req.ID)
		if err != nil {
			uc.log.Error().Msgf("failed connect to gRPC server: %v", err)
			return entity.ErrorGUIResponse(err, "server_id", req.ID)
		}

		uc.addInfoMessage(&entity.Info{Message: entity.MsgServerReflectionInfo})

		uc.services, err = uc.grpcClient.LoadFromReflection()
		if err != nil {
			uc.clearInfoMessages()
			uc.log.Error().Msgf("failed to get services: %v", err)
			return entity.ErrorGUIResponse(err)
		}

		uc.clearInfoMessages()
	} else {
		var protoErr *entity.ProtobufError
		uc.grpcClient.AddProtobuf(uc.curServer.ProtoFiles...)
		uc.grpcClient.AddImport(uc.curServer.ImportPath...)
		uc.services, warn, protoErr = uc.grpcClient.LoadFromProtobuf()
		if protoErr != nil {
			uc.log.Error().Msgf("failed to get services: %v", protoErr.Err)
			return &entity.GUIResponse{
				Status: entity.GUIResponseStatusError,
				Error: entity.Error{
					Code:            protoErr.Code,
					CodeDescription: protoErr.CodeDescription,
					Pos:             protoErr.Pos,
					Message:         protoErr.Error(),
				},
			}
		}
	}

	server.Breadcrumb, err = workspaceUseCase.GetBreadcrumb(req.ID)
	if err != nil {
		uc.log.Error().Msgf("failed to make breadcrumb: %v", err)
		return entity.ErrorGUIResponse(err)
	}

	return &entity.GUIResponse{
		Status: entity.GUIResponseStatusOK,
		Payload: &entity.LoadServerResponse{
			Server:   server,
			Services: uc.services,
			Query:    query,
			Warning:  warn,
		},
	}
}

func (uc *GrpcUseCase) Query(payload map[string]interface{}) *entity.GUIResponse {
	req := &entity.Query{}
	if err := req.Model(payload); err != nil {
		return entity.ErrorGUIResponse(err)
	}

	if err := uc.connect(req.ServerID); err != nil {
		uc.log.Error().Msgf("failed connect to gRPC server: %v", err)
		return entity.ErrorGUIResponse(err)
	}

	uc.clearInfoMessages()

	method, err := uc.getMethodByName(req.Service, req.Method)
	if err != nil {
		return entity.ErrorGUIResponse(err)
	}

	uc.initQuery(req)
	err = uc.grpcClient.Query(method, req.Data, req.Metadata)
	if errors.Is(err, entity.ErrNotConnected) {
		uc.curConnectedServerID = 0
		uc.clearInfoMessages()
		return uc.Query(payload)
	}

	return &entity.GUIResponse{
		Status: entity.GUIResponseStatusOK,
		Payload: &entity.QueryResponse{
			Sent: uc.grpcClient.GetSentCounter(),
		},
	}
}

func (uc *GrpcUseCase) connect(serverID int64) error {
	if uc.curServer.IsK8SEnabled() {
		createForward := true
		existsForward := uc.getPortForward(uc.curServer)
		if existsForward != nil {
			if existsForward.hash != getPortForwardHash(uc.curServer) {
				existsForward.control.Close()
			} else {
				createForward = false
			}
		}

		if createForward {
			uc.curConnectedServerID = 0
			uc.addInfoMessage(&entity.Info{Message: entity.MsgCreatingPortForward})
			uc.curServer.K8SPortForward.ErrHandler = uc.getPortForwardErrorHandler(*uc.curServer, serverID)
			control, err := uc.k8sClient.PortForward(uc.curServer.K8SPortForward)
			if err != nil {
				uc.clearInfoMessages()
				uc.log.Error().Msgf("failed to create k8s port forward: %v", err)
				return err
			}
			uc.addPortForward(uc.curServer, control)
		}
	}

	if uc.curConnectedServerID == serverID {
		return nil
	}

	uc.grpcClient.Close()

	uc.addInfoMessage(&entity.Info{Message: entity.MsgConnectingServer})

	err := uc.grpcClient.Connect(uc.curServer.Addr, uc.curServer.Auth, uc.curServerClientOptions...)
	if err != nil {
		uc.clearInfoMessages()
		uc.log.Error().Msgf("failed connect to gRPC server: %v", err)
		return err
	}

	uc.curConnectedServerID = serverID

	return nil
}

// CancelQuery aborting a running request
func (uc *GrpcUseCase) CancelQuery() {
	uc.grpcClient.CancelQuery()
}

// CloseStream stops a running gRPC stream
func (uc *GrpcUseCase) CloseStream() {
	uc.grpcClient.CloseStream()
}

func (uc *GrpcUseCase) initQuery(q *entity.Query) {
	if uc.curServerID == q.ServerID && uc.curService == q.Service && uc.curMethod == q.Method {
		return
	}

	uc.curServerID = q.ServerID
	uc.curService = q.Service
	uc.curMethod = q.Method

	uc.CancelQuery()
}

func (uc *GrpcUseCase) getServiceByName(serviceName string) (*entity.Service, error) {
	if uc.services == nil {
		return nil, errors.New("services not initialized")
	}

	for _, s := range uc.services {
		if s.Name == serviceName {
			return s, nil
		}
	}

	return nil, fmt.Errorf("service \"%s\" not found", serviceName)
}

func (uc *GrpcUseCase) getMethodByName(serviceName, methodName string) (*entity.Method, error) {
	service, err := uc.getServiceByName(serviceName)
	if err != nil {
		return nil, err
	}

	for _, m := range service.Methods {
		if m.Name == methodName {
			return m, nil
		}
	}

	return nil, fmt.Errorf("method \"%s.%s\" not found", serviceName, methodName)
}

func (uc *GrpcUseCase) getPortForwardErrorHandler(srv entity.WorkspaceItemServer, serverID int64) func(error) {
	return func(err error) {
		uc.deletePortForward(srv)
		if uc.curConnectedServerID == serverID {
			uc.curConnectedServerID = 0
			uc.errorCh <- &entity.Error{
				Message: err.Error(),
			}
		}
	}
}

func (uc *GrpcUseCase) getPortForward(srv *entity.WorkspaceItemServer) *forwardPort {
	uc.muForwardPorts.RLock()
	defer uc.muForwardPorts.RUnlock()

	if uc.forwardPorts == nil {
		return nil
	}

	if fp, ok := uc.forwardPorts[srv.K8SPortForward.LocalPort]; ok {
		return fp
	}

	return nil
}

func (uc *GrpcUseCase) addPortForward(srv *entity.WorkspaceItemServer, control entity.PortForwardControl) {
	uc.muForwardPorts.Lock()
	defer uc.muForwardPorts.Unlock()

	if uc.forwardPorts == nil {
		uc.forwardPorts = make(map[int16]*forwardPort, 10)
	}

	uc.forwardPorts[srv.K8SPortForward.LocalPort] = &forwardPort{
		control: control,
		hash:    getPortForwardHash(srv),
	}
}

func (uc *GrpcUseCase) deletePortForward(srv entity.WorkspaceItemServer) {
	uc.muForwardPorts.Lock()
	defer uc.muForwardPorts.Unlock()

	if uc.forwardPorts == nil {
		return
	}

	delete(uc.forwardPorts, srv.K8SPortForward.LocalPort)
}

// GetInfoChannel returns info channel
func (uc *GrpcUseCase) GetInfoChannel() chan *entity.Info {
	return uc.infoCh
}

// GetErrorChannel returns error channel
func (uc *GrpcUseCase) GetErrorChannel() chan *entity.Error {
	return uc.errorCh
}

func (uc *GrpcUseCase) addInfoMessage(m *entity.Info) {
	uc.infoCh <- m
}

func (uc *GrpcUseCase) clearInfoMessages() {
	uc.infoCh <- &entity.Info{}
}

func getPortForwardHash(srv *entity.WorkspaceItemServer) string {
	data := fmt.Sprintf("%d|%s|%s|%s",
		srv.K8SPortForward.PodPort,
		srv.K8SPortForward.Namespace,
		srv.K8SPortForward.PodName,
		srv.K8SPortForward.PodNameSelector,
	)

	if srv.K8SPortForward.ClientConfig.GCSAuth != nil && srv.K8SPortForward.ClientConfig.GCSAuth.Enabled {
		data = fmt.Sprintf("|%s|%s|%s|%s",
			data,
			srv.K8SPortForward.ClientConfig.GCSAuth.Project,
			srv.K8SPortForward.ClientConfig.GCSAuth.Location,
			srv.K8SPortForward.ClientConfig.GCSAuth.Cluster,
		)
	}

	hash := md5.Sum([]byte(data))
	return hex.EncodeToString(hash[:])
}

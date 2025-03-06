// Package usecase provides business logic.
package usecase

import (
	"context"
	"fmt"
	"sync"

	"github.com/pkg/errors"

	"github.com/forest33/warthog/adapter/grpc"
	"github.com/forest33/warthog/business/entity"
	"github.com/forest33/warthog/pkg/logger"
)

// GrpcUseCase object capable of interacting with GrpcUseCase.
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
	forwardPorts           map[uint16]*forwardPort
	muForwardPorts         sync.RWMutex
	infoCh                 chan *entity.Info
	errorCh                chan *entity.Error
}

// GrpcClient is an interface for working with the gRPC client.
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

// K8SClient is an interface for working with the K8S.
type K8SClient interface {
	PortForward(r *entity.K8SPortForward) (entity.PortForwardControl, error)
}

type forwardPort struct {
	control entity.PortForwardControl
	hash    string
}

// NewGrpcUseCase creates a new GrpcUseCase.
func NewGrpcUseCase(ctx context.Context, log *logger.Zerolog, grpcClient GrpcClient, k8sClient K8SClient, workspaceRepo WorkspaceRepo) *GrpcUseCase {
	useCase := &GrpcUseCase{
		ctx:           ctx,
		log:           log,
		grpcClient:    grpcClient,
		k8sClient:     k8sClient,
		workspaceRepo: workspaceRepo,
		infoCh:        make(chan *entity.Info),
		errorCh:       make(chan *entity.Error),
	}

	useCase.initSubscriptions()

	return useCase
}

func (uc *GrpcUseCase) initSubscriptions() {
	workspaceUseCase.Subscribe(func(e entity.WorkspaceEvent, payload interface{}) {
		switch e {
		case entity.WorkspaceEventServerUpdated:
			w := payload.(*entity.Workspace)
			if w.ID == uc.curServerID {
				uc.curServer = w.Data.(*entity.WorkspaceItemServer)
				uc.curConnectedServerID = 0
			}
			uc.deletePortForward(*w.Data.(*entity.WorkspaceItemServer))
		default:
			uc.log.Error().Msgf("unknown workspace event: %s", e.String())
		}
	})
}

// LoadServer reads the server description from the database and returns it to the GUI.
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

	uc.curServerID = server.ID
	uc.curServer = server.Data.(*entity.WorkspaceItemServer)
	uc.curServerClientOptions = nil
	uc.curServerClientOptions = make([]grpc.ClientOpt, 0, 4)

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

func (uc *GrpcUseCase) connect(serverID int64) error {
	if uc.curServer.IsK8SEnabled() {
		createForward := true
		existsForward := uc.getPortForward(uc.curServer)
		if existsForward != nil {
			if existsForward.hash != uc.curServer.PortForwardHash() {
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

// GetInfoChannel returns info channel.
func (uc *GrpcUseCase) GetInfoChannel() chan *entity.Info {
	return uc.infoCh
}

// GetErrorChannel returns error channel.
func (uc *GrpcUseCase) GetErrorChannel() chan *entity.Error {
	return uc.errorCh
}

func (uc *GrpcUseCase) addInfoMessage(m *entity.Info) {
	uc.infoCh <- m
}

func (uc *GrpcUseCase) clearInfoMessages() {
	uc.infoCh <- &entity.Info{}
}

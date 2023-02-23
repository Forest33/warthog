// Package usecase provides business logic.
package usecase

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	"github.com/forest33/warthog/adapter/grpc"
	"github.com/forest33/warthog/business/entity"
	"github.com/forest33/warthog/pkg/logger"
)

// GrpcUseCase object capable of interacting with GrpcUseCase
type GrpcUseCase struct {
	ctx                    context.Context
	log                    *logger.Zerolog
	client                 GrpcClient
	services               []*entity.Service
	workspaceRepo          WorkspaceRepo
	curServerID            int64
	curConnectedServerID   int64
	curService             string
	curMethod              string
	curServer              *entity.WorkspaceItemServer
	curServerClientOptions []grpc.ClientOpt
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
	Query(method *entity.Method, data map[string]interface{}, metadata []string)
	CancelQuery()
	CloseStream()
	Close()
}

// NewGrpcUseCase creates a new GrpcUseCase
func NewGrpcUseCase(ctx context.Context, log *logger.Zerolog, client GrpcClient, workspaceRepo WorkspaceRepo) *GrpcUseCase {
	return &GrpcUseCase{
		ctx:           ctx,
		log:           log,
		client:        client,
		workspaceRepo: workspaceRepo,
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
		err = uc.client.Connect(uc.curServer.Addr, uc.curServer.Auth, uc.curServerClientOptions...)
		if err != nil {
			uc.log.Error().Msgf("failed connect to gRPC server: %v", err)
			return entity.ErrorGUIResponse(err)
		}
		uc.curConnectedServerID = req.ID

		uc.services, err = uc.client.LoadFromReflection()
		if err != nil {
			uc.log.Error().Msgf("failed to get services: %v", err)
			return entity.ErrorGUIResponse(err)
		}
	} else {
		var protoErr *entity.ProtobufError
		uc.client.AddProtobuf(uc.curServer.ProtoFiles...)
		uc.client.AddImport(uc.curServer.ImportPath...)
		uc.services, warn, protoErr = uc.client.LoadFromProtobuf()
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

	if uc.curConnectedServerID != req.ServerID {
		err := uc.client.Connect(uc.curServer.Addr, uc.curServer.Auth, uc.curServerClientOptions...)
		if err != nil {
			uc.log.Error().Msgf("failed connect to gRPC server: %v", err)
			return entity.ErrorGUIResponse(err)
		}
		uc.curConnectedServerID = req.ServerID
	}

	method, err := uc.getMethodByName(req.Service, req.Method)
	if err != nil {
		return entity.ErrorGUIResponse(err)
	}

	uc.initQuery(req)
	uc.client.Query(method, req.Data, req.Metadata)

	return &entity.GUIResponse{
		Status: entity.GUIResponseStatusOK,
		Payload: &entity.QueryResponse{
			Sent: uc.client.GetSentCounter(),
		},
	}
}

// CancelQuery aborting a running request
func (uc *GrpcUseCase) CancelQuery() {
	uc.client.CancelQuery()
}

// CloseStream stops a running gRPC stream
func (uc *GrpcUseCase) CloseStream() {
	uc.client.CloseStream()
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

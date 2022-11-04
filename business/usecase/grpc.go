package usecase

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	"warthog/adapter/grpc"
	"warthog/business/entity"
	"warthog/pkg/logger"
)

type GrpcUseCase struct {
	ctx           context.Context
	log           *logger.Zerolog
	client        GrpcClient
	services      []*entity.Service
	workspaceRepo WorkspaceRepo
}

type GrpcClient interface {
	Connect(addr string, opts ...grpc.ClientOpt) error
	AddProtobuf(path ...string)
	AddImport(path ...string)
	LoadFromProtobuf() ([]*entity.Service, error)
	LoadFromReflection() ([]*entity.Service, error)
	Query(method *entity.Method, data map[string]interface{}, metadata []string) (*entity.QueryResponse, error)
	CancelQuery()
	Close()
}

func NewGrpcUseCase(ctx context.Context, log *logger.Zerolog, client GrpcClient, workspaceRepo WorkspaceRepo) *GrpcUseCase {
	return &GrpcUseCase{
		ctx:           ctx,
		log:           log,
		client:        client,
		workspaceRepo: workspaceRepo,
	}
}

func (uc *GrpcUseCase) LoadServer(payload map[string]interface{}) *entity.GUIResponse {
	req := &entity.ServerRequest{}
	if err := req.Model(payload); err != nil {
		return entity.ErrorGUIResponse(err)
	}

	var (
		query  *entity.Workspace
		server *entity.Workspace
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

	serverData := server.Data.(*entity.WorkspaceItemServer)

	opts := make([]grpc.ClientOpt, 0, 4)
	if serverData.NoTLS {
		opts = append(opts, grpc.WithNoTLS())
	} else {
		if serverData.Insecure {
			opts = append(opts, grpc.WithInsecure())
		}
		opts = append(opts,
			grpc.WithRootCertificate(serverData.RootCertificate),
			grpc.WithClientCertificate(serverData.ClientCertificate),
			grpc.WithClientKey(serverData.ClientKey),
		)
	}

	err = uc.client.Connect(serverData.Addr, opts...)
	if err != nil {
		uc.log.Error().Msgf("failed connect to gRPC server: %v", err)
		return entity.ErrorGUIResponse(err)
	}

	if serverData.UseReflection {
		uc.services, err = uc.client.LoadFromReflection()
	} else {
		uc.client.AddProtobuf(serverData.ProtoFiles...)
		uc.client.AddImport(serverData.ImportPath...)
		uc.services, err = uc.client.LoadFromProtobuf()
	}

	if err != nil {
		uc.log.Error().Msgf("failed to get services: %v", err)
		return entity.ErrorGUIResponse(err)
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
		},
	}
}

func (uc *GrpcUseCase) Query(payload map[string]interface{}) *entity.GUIResponse {
	req := &entity.Query{}
	if err := req.Model(payload); err != nil {
		return entity.ErrorGUIResponse(err)
	}

	method, err := uc.getMethodByName(req.Service, req.Method)
	if err != nil {
		return entity.ErrorGUIResponse(err)
	}

	resp, err := uc.client.Query(method, req.Data, req.Metadata)
	if err != nil {
		uc.log.Error().Msgf("failed to execute query: %v", err)
		return entity.ErrorGUIResponse(err)
	}

	return &entity.GUIResponse{
		Status:  entity.GUIResponseStatusOK,
		Payload: resp,
	}
}

func (uc *GrpcUseCase) CancelQuery() {
	uc.client.CancelQuery()
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

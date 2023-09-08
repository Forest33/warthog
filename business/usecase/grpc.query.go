package usecase

import (
	"errors"

	"github.com/forest33/warthog/business/entity"
)

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

	uc.curService = q.Service
	uc.curMethod = q.Method

	uc.CancelQuery()
}

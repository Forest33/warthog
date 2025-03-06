package usecase

import (
	"context"

	"k8s.io/apimachinery/pkg/util/runtime"

	"github.com/forest33/warthog/business/entity"
)

func (uc *GrpcUseCase) getPortForwardErrorHandler(srv entity.WorkspaceItemServer, serverID int64) runtime.ErrorHandler {
	return func(ctx context.Context, err error, msg string, keysAndValues ...interface{}) {
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
		uc.forwardPorts = make(map[uint16]*forwardPort, 10)
	}

	uc.forwardPorts[srv.K8SPortForward.LocalPort] = &forwardPort{
		control: control,
		hash:    srv.PortForwardHash(),
	}
}

func (uc *GrpcUseCase) deletePortForward(srv entity.WorkspaceItemServer) {
	uc.muForwardPorts.Lock()
	defer uc.muForwardPorts.Unlock()

	if uc.forwardPorts == nil {
		return
	}

	if fp, ok := uc.forwardPorts[srv.K8SPortForward.LocalPort]; ok && fp.control != nil {
		fp.control.Close()
	}

	delete(uc.forwardPorts, srv.K8SPortForward.LocalPort)
}

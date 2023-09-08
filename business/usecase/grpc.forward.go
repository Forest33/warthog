package usecase

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"

	"github.com/forest33/warthog/business/entity"
)

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

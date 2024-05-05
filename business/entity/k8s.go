package entity

import (
	"bytes"
	"errors"
	"strconv"
)

// K8SClientConfig k8s client config.
type K8SClientConfig struct {
	// GCSAuth GCP authentication request.
	GCSAuth *GCSAuth `json:"auth"`
	// KubeConfigFile absolute path to the kubernetes config file.
	KubeConfigFile string `json:"config_file"`
	// BearerToken Bearer Token string.
	BearerToken string `json:"bearer_token"`
}

// K8SPortForward k8s port forward request.
type K8SPortForward struct {
	// Enabled k8s port forward.
	Enabled bool `json:"enabled"`
	// k8s client config.
	ClientConfig *K8SClientConfig `json:"client_config"`
	// Namespace is the pod namespace.
	Namespace string `json:"namespace"`
	// PodName is the pod name.
	PodName string `json:"pod_name"`
	// PodNameSelector is the pod name selector.
	PodNameSelector string `json:"pod_name_selector"`
	// LocalPort is the local port that will be selected to expose the PodPort.
	LocalPort uint16 `json:"local_port"`
	// PodPort is the target port for the pod.
	PodPort uint16 `json:"pod_port"`
	// ErrHandler error handler.
	ErrHandler func(err error) `json:"-"`
}

// GCSAuth GCS authentication request.
type GCSAuth struct {
	// Enabled GCS authentication enabled.
	Enabled bool `json:"enabled"`
	// Project GKE project.
	Project string `json:"project"`
	// Location cluster location.
	Location string `json:"location"`
	// Cluster the name of the cluster.
	Cluster string `json:"cluster"`
}

type PortForwardControl interface {
	Close()
	Output() *bytes.Buffer
	Error() *bytes.Buffer
}

// Model creates K8SPortForward from UI request.
func (p *K8SPortForward) Model(req map[string]interface{}) error {
	if req == nil {
		return errors.New("no data")
	}

	if v, ok := req["client_config"]; ok && v != nil {
		if cf, ok := v.(map[string]interface{}); !ok {
			return errors.New("client config not a map[string]interface{}")
		} else {
			p.ClientConfig = &K8SClientConfig{}
			if err := p.ClientConfig.Model(cf); err != nil {
				return err
			}
		}
	} else {
		return errors.New("no client config")
	}

	if v, ok := req["enabled"]; ok && v != nil {
		if p.Enabled, ok = v.(bool); !ok {
			return errors.New("enabled not a boolean")
		}
	}
	if v, ok := req["namespace"]; ok && v != nil {
		if p.Namespace, ok = v.(string); !ok {
			return errors.New("namespace not a string")
		}
	}
	if v, ok := req["pod_name"]; ok && v != nil {
		if p.PodName, ok = v.(string); !ok {
			return errors.New("pod name not a string")
		}
	}
	if v, ok := req["pod_name_selector"]; ok && v != nil {
		if p.PodNameSelector, ok = v.(string); !ok {
			return errors.New("pod name selector not a string")
		}
	}
	if v, ok := req["local_port"]; ok && v != nil {
		if lp, ok := v.(string); !ok {
			return errors.New("local port not a string")
		} else {
			port, err := strconv.ParseInt(lp, 10, 32)
			if err != nil {
				return err
			}
			p.LocalPort = uint16(port)
		}
	}
	if v, ok := req["pod_port"]; ok && v != nil {
		if pp, ok := v.(string); !ok {
			return errors.New("pod port not a string")
		} else {
			port, err := strconv.ParseInt(pp, 10, 32)
			if err != nil {
				return err
			}
			p.PodPort = uint16(port)
		}
	}

	return nil
}

// Model creates GCSAuth from UI request.
func (a *GCSAuth) Model(auth map[string]interface{}) error {
	if auth == nil {
		return errors.New("no data")
	}

	if v, ok := auth["enabled"]; ok && v != nil {
		if a.Enabled, ok = v.(bool); !ok {
			return errors.New("enabled not a boolean")
		}
	}
	if v, ok := auth["project"]; ok && v != nil {
		if a.Project, ok = v.(string); !ok {
			return errors.New("project not a string")
		}
	}
	if v, ok := auth["location"]; ok && v != nil {
		if a.Location, ok = v.(string); !ok {
			return errors.New("location not a string")
		}
	}
	if v, ok := auth["cluster"]; ok && v != nil {
		if a.Cluster, ok = v.(string); !ok {
			return errors.New("cluster not a string")
		}
	}

	return nil
}

// Model creates K8SClientConfig from UI request.
func (c *K8SClientConfig) Model(cfg map[string]interface{}) error {
	if cfg == nil {
		return errors.New("no data")
	}

	if v, ok := cfg["auth"]; ok && v != nil {
		if a, ok := v.(map[string]interface{}); !ok {
			return errors.New("auth not a map[string]interface{}")
		} else {
			c.GCSAuth = &GCSAuth{}
			if err := c.GCSAuth.Model(a); err != nil {
				return err
			}
		}
	}

	if v, ok := cfg["config_file"]; ok && v != nil {
		if c.KubeConfigFile, ok = v.(string); !ok {
			return errors.New("config file not a string")
		}
	}
	if v, ok := cfg["bearer_token"]; ok && v != nil {
		if c.BearerToken, ok = v.(string); !ok {
			return errors.New("bearer token not a string")
		}
	}

	return nil
}

package entity

import (
	"bytes"
	"errors"
	"strconv"
)

// K8SClientConfig k8s client config
type K8SClientConfig struct {
	// GCSAuth GCP authentication request
	GCSAuth *GCSAuth `json:"auth"`
	// KubeConfigFile absolute path to the kubernetes config file
	KubeConfigFile string `json:"config_file"`
	// BearerToken Bearer Token string
	BearerToken string `json:"bearer_token"`
}

// K8SPortForward k8s port forward request
type K8SPortForward struct {
	// Enabled k8s port forward
	Enabled bool `json:"enabled"`
	// k8s client config
	ClientConfig *K8SClientConfig `json:"client_config"`
	// Namespace is the pod namespace
	Namespace string `json:"namespace"`
	// PodName is the pod name
	PodName string `json:"pod_name"`
	// PodNameSelector is the pod name selector
	PodNameSelector string `json:"pod_name_selector"`
	// LocalPort is the local port that will be selected to expose the PodPort
	LocalPort int16 `json:"local_port"`
	// PodPort is the target port for the pod
	PodPort int16 `json:"pod_port"`
	// ErrHandler error handler
	ErrHandler func(err error) `json:"-"`
}

// GCSAuth GCS authentication request
type GCSAuth struct {
	// Enabled GCS authentication enabled
	Enabled bool `json:"enabled"`
	// Project GKE project
	Project string `json:"project"`
	// Location cluster location
	Location string `json:"location"`
	// Cluster the name of the cluster
	Cluster string `json:"cluster"`
}

type PortForwardControl interface {
	Close()
	Output() *bytes.Buffer
	Error() *bytes.Buffer
}

// Model creates K8SPortForward from UI request
func (p *K8SPortForward) Model(req map[string]interface{}) error {
	if req == nil {
		return errors.New("no data")
	}

	if v, ok := req["client_config"]; ok && v != nil {
		p.ClientConfig = &K8SClientConfig{}
		if err := p.ClientConfig.Model(v.(map[string]interface{})); err != nil {
			return err
		}
	} else {
		return errors.New("no client config")
	}

	if v, ok := req["enabled"]; ok && v != nil {
		p.Enabled = v.(bool)
	}
	if v, ok := req["namespace"]; ok && v != nil {
		p.Namespace = v.(string)
	}
	if v, ok := req["pod_name"]; ok && v != nil {
		p.PodName = v.(string)
	}
	if v, ok := req["pod_name_selector"]; ok && v != nil {
		p.PodNameSelector = v.(string)
	}
	if v, ok := req["local_port"]; ok && v != nil {
		port, err := strconv.ParseInt(v.(string), 10, 32)
		if err != nil {
			return err
		}
		p.LocalPort = int16(port)
	}
	if v, ok := req["pod_port"]; ok && v != nil {
		port, err := strconv.ParseInt(v.(string), 10, 32)
		if err != nil {
			return err
		}
		p.PodPort = int16(port)
	}

	return nil
}

// Model creates GCSAuth from UI request
func (a *GCSAuth) Model(auth map[string]interface{}) error {
	if auth == nil {
		return errors.New("no data")
	}

	if v, ok := auth["enabled"]; ok && v != nil {
		a.Enabled = v.(bool)
	}
	if v, ok := auth["project"]; ok && v != nil {
		a.Project = v.(string)
	}
	if v, ok := auth["location"]; ok && v != nil {
		a.Location = v.(string)
	}
	if v, ok := auth["cluster"]; ok && v != nil {
		a.Cluster = v.(string)
	}

	return nil
}

// Model creates K8SClientConfig from UI request
func (c *K8SClientConfig) Model(cfg map[string]interface{}) error {
	if cfg == nil {
		return errors.New("no data")
	}

	if v, ok := cfg["auth"]; ok && v != nil {
		c.GCSAuth = &GCSAuth{}
		if err := c.GCSAuth.Model(v.(map[string]interface{})); err != nil {
			return err
		}
	}

	if v, ok := cfg["config_file"]; ok && v != nil {
		c.KubeConfigFile = v.(string)
	}
	if v, ok := cfg["bearer_token"]; ok && v != nil {
		c.BearerToken = v.(string)
	}

	return nil
}

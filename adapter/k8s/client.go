package k8s

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/portforward"
	"k8s.io/client-go/transport/spdy"

	"github.com/forest33/warthog/business/entity"
	"github.com/forest33/warthog/pkg/logger"
)

// Client object capable of interacting with Client.
type Client struct {
	ctx context.Context
	cfg *entity.Settings
	log *logger.Zerolog
}

// New creates a new Client.
func New(ctx context.Context, log *logger.Zerolog) *Client {
	return &Client{
		ctx: ctx,
		log: log,
	}
}

// SetSettings sets application settings.
func (c *Client) SetSettings(cfg *entity.Settings) {
	c.cfg = cfg
}

// PortForward port forward.
func (c *Client) PortForward(r *entity.K8SPortForward) (entity.PortForwardControl, error) {
	config, client, err := c.createClient(r.ClientConfig)
	if err != nil {
		return nil, err
	}

	var (
		podName string
	)

	if r.PodName != "" {
		podName = r.PodName
	} else if r.PodNameSelector != "" {
		if podName, err = c.findPod(client, r.Namespace, r.PodNameSelector); err != nil {
			return nil, err
		}
	} else {
		return nil, entity.ErrK8SPodNotFound
	}

	ctrl := &PortForwardControl{
		stopCh: make(chan struct{}, 1),
		out:    &bytes.Buffer{},
		errOut: &bytes.Buffer{},
	}

	readyCh := make(chan struct{})

	writeError := func(err error) {
		if _, err := ctrl.errOut.Write([]byte(err.Error())); err != nil {
			c.log.Error().Msgf("failed write to error stream: %v", err)
		}
		close(readyCh)
	}

	if r.ErrHandler != nil {
		runtime.ErrorHandlers = []func(error){r.ErrHandler}
	}

	go func() {
		path := fmt.Sprintf("/api/v1/namespaces/%s/pods/%s/portforward", r.Namespace, podName)

		u, err := url.Parse(config.Host)
		if err != nil {
			writeError(err)
			return
		}

		transport, upgrader, err := spdy.RoundTripperFor(config)
		if err != nil {
			writeError(err)
			return
		}

		httpClient := &http.Client{
			Transport: transport,
			Timeout:   time.Duration(*c.cfg.K8SRequestTimeout) * time.Second,
		}

		dialer := spdy.NewDialer(upgrader, httpClient, http.MethodPost, &url.URL{Scheme: "https", Path: path, Host: u.Host})
		fw, err := portforward.New(dialer, []string{fmt.Sprintf("%d:%d", r.LocalPort, r.PodPort)}, ctrl.stopCh, readyCh, ctrl.out, ctrl.errOut)
		if err != nil {
			writeError(err)
			return
		}

		if err := fw.ForwardPorts(); err != nil {
			writeError(err)
			return
		}
	}()

	<-readyCh

	if ctrl.errOut.Len() != 0 {
		return nil, errors.New(ctrl.errOut.String())
	}

	return ctrl, nil
}

func (c *Client) createClient(cfg *entity.K8SClientConfig) (*rest.Config, *kubernetes.Clientset, error) {
	var (
		restConfig *rest.Config
		err        error
	)

	if cfg.GCSAuth != nil && cfg.GCSAuth.Enabled {
		ctx, cancel := context.WithTimeout(c.ctx, time.Duration(*c.cfg.K8SRequestTimeout)*time.Second)
		defer cancel()

		restConfig, err = c.gcsAuth(ctx, cfg.GCSAuth)
		if err != nil {
			return nil, nil, err
		}
	} else {
		kubeConfig := cfg.KubeConfigFile
		if kubeConfig == "" {
			if home := homeDir(); home != "" {
				kubeConfig = filepath.Join(home, ".kube", "config")
			}
		}

		restConfig, err = clientcmd.BuildConfigFromFlags("", kubeConfig)
		if err != nil {
			return nil, nil, err
		}

		restConfig.BearerToken = cfg.BearerToken
	}

	restConfig.Timeout = time.Duration(*c.cfg.K8SRequestTimeout) * time.Second

	client, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return nil, nil, err
	}

	return restConfig, client, nil
}

func (c *Client) findPod(client *kubernetes.Clientset, namespace, selector string) (string, error) {
	ctx, cancel := context.WithTimeout(c.ctx, time.Duration(*c.cfg.K8SRequestTimeout)*time.Second)
	defer cancel()

	pods, err := client.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{LabelSelector: selector})
	if err != nil {
		return "", err
	}

	if len(pods.Items) == 0 {
		return "", entity.ErrK8SPodNotFound
	}

	return pods.Items[0].Name, nil
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE")
}

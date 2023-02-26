package k8s

import (
	"bytes"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type PortForwardControl struct {
	config *rest.Config
	client *kubernetes.Clientset
	stopCh chan struct{}
	out    *bytes.Buffer
	errOut *bytes.Buffer
}

func (c *PortForwardControl) Close() {
	close(c.stopCh)
}

func (c *PortForwardControl) Output() *bytes.Buffer {
	return c.out
}

func (c *PortForwardControl) Error() *bytes.Buffer {
	return c.errOut
}

package k8s

import (
	"bytes"
)

// PortForwardControl port forwarding control.
type PortForwardControl struct {
	stopCh chan struct{}
	out    *bytes.Buffer
	errOut *bytes.Buffer
}

// Close stop port forwarding.
func (c *PortForwardControl) Close() {
	close(c.stopCh)
}

// Output returns a buffer containing port forward messages.
func (c *PortForwardControl) Output() *bytes.Buffer {
	return c.out
}

// Error returns a buffer containing error messages for port forwarding.
func (c *PortForwardControl) Error() *bytes.Buffer {
	return c.errOut
}

package entity

import "errors"

var (
	// ErrK8SPodNotFound error - pod not found.
	ErrK8SPodNotFound = errors.New("pod not found")
	// ErrNotConnected error - server not connected.
	ErrNotConnected = errors.New("not connected")
)

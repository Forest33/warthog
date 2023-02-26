package entity

import "errors"

var (
	ErrK8SPodNotFound = errors.New("pod not found")
	ErrNotConnected   = errors.New("not connected")
)

package entity

import "errors"

var (
	// ErrK8SPodNotFound error - pod not found.
	ErrK8SPodNotFound = errors.New("pod not found")
	// ErrNotConnected error - server not connected.
	ErrNotConnected = errors.New("not connected")
	// ErrFolderAlreadyExists - folder with the same name already exists
	ErrFolderAlreadyExists = errors.New("A folder with the same name already exists")
)

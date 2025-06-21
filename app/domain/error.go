package domain

import (
	"errors"
)

var (
	// ErrNotFound will throw if the requested item is not exists
	ErrNotFound = errors.New("your requested item is not found")
	// ErrAlreadyExists will throw if the requested item is already exists
	ErrAlreadyExists = errors.New("your requested item is already exists")
	// ErrInvalidParameters will throw if the given request-body or params is not valid
	ErrInvalidParameters = errors.New("parameters is not valid")
)

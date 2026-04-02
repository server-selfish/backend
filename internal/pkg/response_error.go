package pkg

import "errors"

var (
	ErrBadRequest   = errors.New("invalid request body")
	ErrAlreadyExist = errors.New("resource is already exist")
	ErrNotFound     = errors.New("resource already exist")
)

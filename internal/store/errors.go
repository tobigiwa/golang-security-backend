package store

import (
	"errors"
)

var (
	ErrDuplicateEmail     = errors.New("duplicate email")
	ErrDuplicateUsername  = errors.New("duplicate username")
	ErrTimeOut            = errors.New("timeout")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

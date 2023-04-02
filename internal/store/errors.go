package store

import (
	"errors"
)

var (
	ErrDuplicateEmail                = errors.New("duplicate email")
	ErrDuplicateUsername             = errors.New("duplicate username")
	ErrTimeOut                       = errors.New("timeout")
	ErrInvalidCredentials            = errors.New("invalid credentials")
	ErrIncompleteDatabaseCredentials = errors.New("incomplete database credentials")
	ErrLoadingEnvFile                = errors.New("cannot load .env file")
	ErrInvalidUserCredentials        = errors.New("invalid user credentials")
)

package store

import (
	"errors"
)

var (
	ErrDuplicateEmail     = errors.New("database: duplicate email")
	ErrDuplicateUsername  = errors.New("database:duplicate username")
	ErrTimeOut            = errors.New("timeout")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

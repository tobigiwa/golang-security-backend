package models

import (
	"errors"
)

var (
	ErrDuplicateEmail    = errors.New("database: duplicate email")
	ErrDuplicateUsername = errors.New("database:duplicate username")
	ErrTimeOut           = errors.New("timeout")
)

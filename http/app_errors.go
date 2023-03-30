package http

import (
	"errors"
	"net/http"
)

var (
	errDuplicateEmail     = errors.New("database: duplicate email")
	errDuplicateUsername  = errors.New("database:duplicate username")
	errInvalidCredentials = errors.New("invalid credentials")
)

func (a *WebApp) ClientError(w http.ResponseWriter, httpStatus int, text string) {
	http.Error(w, text, httpStatus)
}

func (a *WebApp) ServerError(w http.ResponseWriter, text string) {
	http.Error(w, text, http.StatusInternalServerError)
}


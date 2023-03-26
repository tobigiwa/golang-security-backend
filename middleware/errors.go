package main

import (
	"errors"
	"net/http"
)

var (
	errDuplicateEmail     = errors.New("database: duplicate email")
	errDuplicateUsername  = errors.New("database:duplicate username")
	errTimeOut            = errors.New("timeout")
	errInvalidCredentials = errors.New("invalid credentials")
)

// clientError denotes the any error from the client side (e.g inputs).
func clientError(w http.ResponseWriter, httpStatus int, text string) {
	http.Error(w, text, httpStatus)
}

// serverError denotes the any error from the server side (e.g expected errors).
func serverError(w http.ResponseWriter, text string) {
	http.Error(w, text, http.StatusInternalServerError)
}

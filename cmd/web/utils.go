package main

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/tobigiwa/golang-security-backend/internal/models"
	"golang.org/x/crypto/bcrypt"
)

func (a *WebApp) clientError(w http.ResponseWriter, httpStatus int, text string) {
	// w.WriteHeader(httpStatus)
	http.Error(w, text, httpStatus)
}

// func (a *WebApp) serverError(w http.ResponseWriter, httpStatus int) {
// 	w.WriteHeader(httpStatus)
// 	http.Error(w, http.StatusText(httpStatus), httpStatus)
// }

func (a *WebApp) undefinedError(w http.ResponseWriter, text string) {
	w.WriteHeader(http.StatusNotImplemented)
	http.Error(w, text, http.StatusNotImplemented)
}

func (a *WebApp) generateHashedPassword(password string) ([]byte, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 15)
	return hashedPassword, err
}


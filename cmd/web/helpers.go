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

func (a *WebApp) authenticate(email, password string) error {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	hashedPassword, err := a.dbModel.FetchUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			return models.ErrInvalidCredentials
		} else {
			return err
		}
	}
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return models.ErrInvalidCredentials
		} else {
			return err
		}
	}
	return nil

}

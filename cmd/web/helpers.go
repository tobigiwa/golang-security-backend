package main

import (
	"net/http"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func (a *WebApp) clientError(w http.ResponseWriter, httpStatus int) {
	w.WriteHeader(httpStatus)
	http.Error(w, http.StatusText(httpStatus), httpStatus)
}

func (a *WebApp) serverError(w http.ResponseWriter, httpStatus int) {
	w.WriteHeader(httpStatus)
	http.Error(w, http.StatusText(httpStatus), httpStatus)
}

func (a *WebApp) undefinedError(w http.ResponseWriter, err string) {
	w.WriteHeader(http.StatusNotImplemented)
	http.Error(w, err, http.StatusNotImplemented)
}

func (a *WebApp) generateHashedPassword(password string) ([]byte, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 15)
	return hashedPassword, err
}

func (a *WebApp) generateNewUUID() string {
	id := uuid.New()
	return id.String()
}

func (a *WebApp) authenticate(email, password string) error {

}

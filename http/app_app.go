package http

import (
	"github.com/tobigiwa/golang-security-backend/internal/store"
	"github.com/tobigiwa/golang-security-backend/logging"
	"golang.org/x/crypto/bcrypt"
)

// Webpp is application struct
type WebApp struct {
	Store  *store.Store
	Logger *logging.Logger
}

func (a *WebApp) generateHashedPassword(password string) ([]byte, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 15)
	return hashedPassword, err
}

package http

import (
	"context"
	"errors"
	"time"

	"github.com/tobigiwa/golang-security-backend/internal/models"
	"golang.org/x/crypto/bcrypt"
)

// Webpp is application struct
type WebApp struct {
	DbModel *models.UserModel
}

func (a *WebApp) generateHashedPassword(password string) ([]byte, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 15)
	return hashedPassword, err
}

func (a *WebApp) Authenticate(email, password string) error {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	hashedPassword, err := a.DbModel.FetchUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, errInvalidCredentials) {
			return errInvalidCredentials
		} else {
			return err
		}
	}
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return errInvalidCredentials
		} else {
			return err
		}
	}
	return nil

}

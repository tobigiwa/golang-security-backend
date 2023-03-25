package middleware

import (
	"context"
	"errors"
	"time"

	"github.com/tobigiwa/golang-security-backend/internal/models"
	"golang.org/x/crypto/bcrypt"
)

// authenticate compares passed password with db password
// returns error if checks fail.
func authenticate(email, password string) error {

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

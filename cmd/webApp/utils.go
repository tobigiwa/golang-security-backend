package main

import (
	"context"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func (a *WebApp) generateHashedPassword(password string) ([]byte, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 15)
	return hashedPassword, err
}

// authenticate compares passed password with db password
// returns error if checks fail.
func (a *WebApp) authenticate(email, password string) error {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	hashedPassword, err := a.dbModel.FetchUserByEmail(ctx, email)
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

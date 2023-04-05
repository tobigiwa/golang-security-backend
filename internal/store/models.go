package store

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type UserModel struct {
	Email    string
	Username string
	Password []byte
	Status   string
}

func (u *UserModel) validate_password(email, password string) error {
	err := bcrypt.CompareHashAndPassword(u.Password, []byte(password))
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return ErrInvalidUserCredentials
	} else {
		return err
	}
}

func (u *UserModel) generateHashedPassword(password string) ([]byte, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 15)
	return hashedPassword, err
}

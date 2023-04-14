package service

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

func (u *UserModel) validatePassword(hashedPassword []byte, password string) error {
	err := bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return ErrInvalidUserCredentials
		default:
			return err
		}
	}
	return nil
}

func (u *UserModel) generateHashedPassword(password string) ([]byte, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 15)
	return hashedPassword, err
}

func (u *UserModel) createUser() string {
	return `INSERT INTO public.user_tbl(email, username, pswd, status) VALUES($1, $2, $3, 'user')`
}

func (u *UserModel) createSuperUser() string {
	return `INSERT INTO public.user_tbl(email, username, pswd, status) VALUES($1, $2, $3, 'superuser')`
}

func (u *UserModel) fetchUser() string {
	return `SELECT email, username, pswd, status FROM public.user_tbl WHERE email = $1`
}

func (u *UserModel) fetchAllUser() string {
	return `SELECT email, username, status FROM public.user_tbl`
}

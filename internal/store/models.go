package store

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type UserModel struct {
	Email    string
	Username string
	Password []byte
	Role     string
}

func (u *UserModel) validateUser(user UserModel, password string) error {
	err := bcrypt.CompareHashAndPassword(u.Password, []byte(password))
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
	stmt := `INSERT INTO public.user_tbl(email, username, pswd, role)
				VALUES($1, $2, $3, 'user')`
	return stmt
}

func (u *UserModel) createSuperUser() string {
	stmt := `INSERT INTO public.user_tbl(email, username, pswd, role)
				VALUES($1, $2, $3, 'superuser')`
	return stmt
}

func (u *UserModel) fetchUser() string {
	stmt := `SELECT email, username, pswd, role FROM public.user_tbl 
	WHERE (email = $1) OR (username = $1)`
	return stmt
}

func (u *UserModel) fetchAllUser() string {
	stmt := `SELECT email, username, role FROM public.user_tbl`
	return stmt
}

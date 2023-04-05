package store

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tobigiwa/golang-security-backend/pkg/logging"
)

type Store struct {
	DB     *pgxpool.Pool
	Logger *logging.Logger
	User *UserModel
}

func (s *Store) InsertUser(email, username, password string) error {
	stmt := `INSERT INTO public.user_tbl(email, username, pswd, status)
				VALUES($1, $2, $3, 'public_user')`
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	hashedPassword, err := s.User.generateHashedPassword(password)
	if err != nil {
		s.Logger.LogError(err, "DB")
		panic(err)
	}
	_, err = s.DB.Exec(ctx, stmt, email, username, hashedPassword)
	if err != nil {
		var pgxError *pgconn.PgError
		if errors.As(err, &pgxError) {
			if pgxError.Code == "23505" && strings.Contains(pgxError.Detail, email) {
				return ErrDuplicateEmail
			} else if pgxError.Code == "23505" && strings.Contains(pgxError.Detail, username) {
				return ErrDuplicateUsername
			} else {
				s.Logger.LogError(err, "DB")
				return err
			}
		}
	}
	return nil
}

func (s *Store) FetchUser(search string) (UserModel, error) {
	stmt := `SELECT email, username, pswd, status FROM public.user_tbl 
				WHERE (email = $1) OR (username = $1)`
	var user UserModel
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	err := s.DB.QueryRow(ctx, stmt, search).Scan(&user.Username, &user.Password, &user.Status)
	defer cancel()
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return UserModel{}, ErrInvalidCredentials
		} else {
			s.Logger.LogError(err, "DB")
			return UserModel{}, err
		}
	}
	return user, nil
}

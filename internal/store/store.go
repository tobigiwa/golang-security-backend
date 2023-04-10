package store

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tobigiwa/golang-security-backend/logging"
)

type Store struct {
	DB     *pgxpool.Pool
	Logger *logging.Logger
	User   *UserModel
}

// PUBLIC API

func (s *Store) CreateSuperUser(email, username, password string) error {
	hashedPassword, err := s.User.generateHashedPassword(password)
	if err != nil {
		s.Logger.LogError(err, "DB")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = s.insertUser(ctx, s.User.createSuperUser(), email, username, string(hashedPassword))
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) CreateUser(email, username, password string) error {
	hashedPassword, err := s.User.generateHashedPassword(password)
	if err != nil {
		s.Logger.LogError(err, "DB")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = s.insertUser(ctx, s.User.createUser(), email, username, string(hashedPassword))
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) FetchUser(search string) (UserModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	user, err := s.fetchUser(ctx, s.User.fetchUser(), search)
	if err != nil {
		return UserModel{}, err
	}
	return user, nil
}

func (s *Store) FetchAllUser() ([]UserModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	users, err := s.fetchAllUser(ctx, s.User.fetchAllUser())
	return users, err
}

// PRIVATE API

func (s *Store) fetchAllUser(ctx context.Context, stmt string) ([]UserModel, error) {
	var list []UserModel
	rows, err := s.DB.Query(ctx, stmt)
	if err != nil {
		return []UserModel{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var r UserModel
		err := rows.Scan(&r.Email, &r.Username, &r.Role)
		if err != nil {
			// do something
			continue
		}
		list = append(list, r)
	}
	if err := rows.Err(); err != nil {
		s.Logger.LogError(err, "DB")
	}
	return list, nil
}

func (s *Store) insertUser(ctx context.Context, stmt, email, username, hashedPassword string) error {
	_, err := s.DB.Exec(ctx, stmt, email, username, hashedPassword)
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

func (s *Store) fetchUser(ctx context.Context, stmt, search string) (UserModel, error) {
	var user UserModel
	err := s.DB.QueryRow(ctx, stmt, search).Scan(&user.Username, &user.Password, &user.Role)
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

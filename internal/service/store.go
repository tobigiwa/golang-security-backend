package service

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

func (s *Store) FetchAllUser() ([]UserModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	users, err := s.fetchAllUser(ctx, s.User.fetchAllUser())
	return users, err
}

func (s *Store) GetAndValidateUser(ctx context.Context, key, password string) (UserModel, error) {
	user, err := s.fetchUser(ctx, s.User.fetchUser(), key)
	if err != nil {
		return user, err
	}
	err = s.User.validatePassword(user.Password, password)
	if err != nil {
		return user, err
	}
	return user, nil
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
		err := rows.Scan(&r.Email, &r.Username, &r.Status)
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
			if pgxError.Code == "23505" {
				switch {
				case strings.Contains(pgxError.Detail, email):
					return ErrDuplicateEmail
				case strings.Contains(pgxError.Detail, username):
					return ErrDuplicateUsername
				}
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
	err := s.DB.QueryRow(ctx, stmt, search).Scan(&user.Email, &user.Username, &user.Password, &user.Status)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return UserModel{}, ErrNotFound
		default:
			s.Logger.LogError(err, "SERVICE")
			return UserModel{}, err
		}
	}
	return user, nil
}

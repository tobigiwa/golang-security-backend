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

type UserModel struct {
	DB     *pgxpool.Pool
	Logger *logging.Logger
}

func New() (*pgxpool.Pool, error) {
	databaseURL, err := dbDSN()
	if err != nil {
		return nil, err
	}
	db, err := pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.Ping(ctx); err != nil {
		return nil, err
	}
	// db.SetMaxOpenConns(20)
	// db.SetMaxIdleConns(10)
	return db, nil
}

func (m *UserModel) Insert(email, username, password string) error {

	stmt := `INSERT INTO public.user_tbl(email, username, pswd, status)
				VALUES($1, $2, $3, 'public_user')`
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := m.DB.Exec(ctx, stmt, email, username, password)
	if err != nil {
		var pgxError *pgconn.PgError
		if errors.As(err, &pgxError) {
			if pgxError.Code == "23505" && strings.Contains(pgxError.Detail, email) {
				return ErrDuplicateEmail
			} else if pgxError.Code == "23505" && strings.Contains(pgxError.Detail, username) {
				return ErrDuplicateUsername
			} else {
				m.Logger.LogError(err, "DB")
				return err
			}
		}
	}
	return nil
}

func (m *UserModel) FetchUserByEmail(ctx context.Context, email string) ([]byte, error) {
	var hashedPassword []byte
	stmt := `SELECT pswd FROM public.user_tbl WHERE email = $1`
	err := m.DB.QueryRow(ctx, stmt, email).Scan(&hashedPassword)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrInvalidCredentials
		} else {
			m.Logger.LogError(err, "DB")
			return nil, err
		}
	}
	return hashedPassword, nil
}

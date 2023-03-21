package models

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserModel struct {
	DB *pgxpool.Pool
}

func (m *UserModel) Insert(uuID, email, username, password string) error {

	stmt := `INSERT INTO public.user_tbl(id, email, username, pswd, status)
				VALUES($1, $2, $3, $4, 'good boy')`
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := m.DB.Exec(ctx, stmt, uuID, email, username, password)
	if err != nil {
		var pgxError *pgconn.PgError
		if errors.As(err, &pgxError) {
			if pgxError.Code == "23505" && strings.Contains(pgxError.Detail, email) {
				return ErrDuplicateEmail
			} else if pgxError.Code == "23505" && strings.Contains(pgxError.Detail, username) {
				return ErrDuplicateUsername
			} else {
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
			return nil, err
		}
	}
	return hashedPassword, nil
}

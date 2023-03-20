package models

import (
	"context"
	"errors"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserModel struct {
	DB *pgxpool.Pool
}

func (m *UserModel) Insert(ctx context.Context, email, username, password string) error {

	stmt := `INSERT INTO public.model_user(email, username, pswd, status)
				VALUES($1, $2, $3, 'good boy')`
	_, err := m.DB.Exec(ctx, stmt, email, username, password)
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

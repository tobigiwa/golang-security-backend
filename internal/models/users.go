package models

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Users struct {
	Email           string `json:"email"`
	Username        string `json:"username"`
	Hashed_password string `json:"passwor"`
	Status          string `json:"status"`
}

type UserModel struct {
	DB *pgxpool.Pool
}

func (m *UserModel) Insert(email, username, password string) error {

	stmt := `INSERT INTO public.model_user(email, username, hashed_password, about)
				VALUES($1, $2, $3, 'good boy')`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := m.DB.Exec(ctx, stmt, email, username, password)

	var pgxError *pgconn.PgError
	if errors.As(err, &pgxError) {
		if pgxError.Code == "23505" && strings.Contains(pgxError.Detail, email) {
			fmt.Println("email")
		} else if pgxError.Code == "23505" && strings.Contains(pgxError.Detail, username) {
			fmt.Println("ussername")

		}
	}

	return err
}

package service

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func NewDatabaseConn() (*pgxpool.Pool, error) {
	databaseURL, err := databseDSN()
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
	return db, nil
}

func databseDSN() (string, error) {
	err := godotenv.Load()
	if err != nil {
		return "", ErrLoadingEnvFile
	}
	var hold = make(map[string]string, 3)
	for _, v := range [3]string{"HOSTNAME", "PASSWORD", "DBNAME"} {
		if val, ok := os.LookupEnv(v); ok && val != "" {
			hold[v] = val
		} else {
			return "", ErrIncompleteDatabaseCredentials
		}
	}
	dbLogin := fmt.Sprintf("postgres://%v:%v@localhost:5432/%v", hold["HOSTNAME"], hold["PASSWORD"], hold["DBNAME"])
	return dbLogin, nil
}

package store

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func DbSetUp() (*pgxpool.Pool, error) {
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

func dbDSN() (string, error) {
	err := godotenv.Load()
	if err != nil {
		return "", errors.New("cannot load .env file")
	}
	var hold = make(map[string]string, 3)
	for _, v := range [3]string{"HOSTNAME", "PASSWORD", "DBNAME"} {
		if val, ok := os.LookupEnv(v); ok && val != "" {
			hold[v] = val
		} else {
			return "", errors.New("incomplete database credentials")
		}
	}

	dbLogin := fmt.Sprintf("postgres://%v:%v@localhost:5432/%v", hold["HOSTNAME"], hold["PASSWORD"], hold["DBNAME"])
	return dbLogin, nil
}

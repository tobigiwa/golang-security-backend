package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func main() {

	db, err := dbSetUp()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Print(db)
}

func dbSetUp() (*sql.DB, error) {

	databaseURL, err := dbDSN()
	if err != nil {
		log.Fatalf("error was: %v", err)
	}

	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		log.Fatalf("unable to establish database connection: %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	return db, nil

}

// dbDSN loads the database credentials from .env file at root directory
// and returns the database connection URL(string) and an error (if any).
func dbDSN() (string, error) {

	err := godotenv.Load()
	if err != nil {
		return "", err
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

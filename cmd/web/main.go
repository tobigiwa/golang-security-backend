package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/tobigiwa/golang-security-backend/internal/models"
)

type WebApp struct {
	dbModel *models.UserModel
	
}

func main() {

	db, err := dbSetUp()
	if err != nil {
		log.Fatal(err)
	}

	app := &WebApp{
		dbModel: &models.UserModel{DB: db},
	}

	webServer := &http.Server{
		Addr:         ":5030",
		Handler:      app.Routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	fmt.Println("server running...")
	err = webServer.ListenAndServe()
	log.Fatal(err)
}

// dbSetUp returns an error or database handle with these settings,
// db.SetMaxOpenConns(20)
// db.SetMaxIdleConns(10).
func dbSetUp() (*pgxpool.Pool, error) {
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

package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	app "github.com/tobigiwa/golang-security-backend/http"
	"github.com/tobigiwa/golang-security-backend/internal/store"
)

func main() {

	db, err := store.DbSetUp()
	if err != nil {
		log.Fatal(err)
	}

	application := &app.WebApp{
		DbModel: &store.UserModel{DB: db},
		
	}

	webServer := &http.Server{
		Addr:         ":5030",
		Handler:      application.Routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	fmt.Println("server running...")
	err = webServer.ListenAndServe()
	log.Fatal(err)
}

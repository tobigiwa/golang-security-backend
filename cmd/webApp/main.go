package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	app "github.com/tobigiwa/golang-security-backend/http"
	"github.com/tobigiwa/golang-security-backend/internal/service"
	"github.com/tobigiwa/golang-security-backend/logging"
)

func main() {

	db, err := service.NewDatabaseConn()
	if err != nil {
		log.Fatal(err)
	}
	logger, err := logging.NewLogger()
	if err != nil {
		var pathError *os.PathError
		if errors.As(err, &pathError) {
			log.Fatal("Incorrect file path for Log")
		}
	}
	application := &app.WebApp{
		Service: &service.Store{DB: db, Logger: logger, User: &service.UserModel{}},
		Logger:  logger,
	}
	webServer := &http.Server{
		Addr:         ":5030",
		Handler:      application.Routes(),
		ErrorLog:     log.New(logger, "", 0),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	fmt.Println("server running...")
	application.Logger.LogInfo("SERVER IS RUNNING", "APP")
	err = webServer.ListenAndServe()
	log.Fatal(err)
}

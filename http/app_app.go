package http

import (
	"encoding/hex"
	"log"
	"net/http"
	"time"

	cookiePackage "github.com/tobigiwa/golang-security-backend/http/cookies"
	"github.com/tobigiwa/golang-security-backend/internal/service"
	"github.com/tobigiwa/golang-security-backend/logging"
)

// Webpp is application struct
type WebApp struct {
	Service *service.Store
	Logger  *logging.Logger
}

func (a *WebApp) CreateCookie(w http.ResponseWriter, payload string) error {
	secretKey, err := hex.DecodeString("13d6b4dff8f84a10851021ec8608f814570d562c92fe6b5ec4c9f595bcb3234b")
	if err != nil {
		log.Fatal(err)
	}

	cookie := http.Cookie{
		Name:    "cookie",
		Value:   payload,
		Expires: time.Now().Add(30 * time.Minute),
		MaxAge:  1800,
		// Secure:   true,
		// HttpOnly: true,
		// SameSite: http.SameSiteLaxMode,
	}
	err = cookiePackage.WriteEncryptCookie(w, cookie, secretKey)
	if err != nil {
		a.ServerError(w, err.Error())
		return err
	}
	return nil
}

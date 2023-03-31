package http

import (
	"encoding/hex"
	"errors"
	"log"
	"net/http"

	"github.com/tobigiwa/golang-security-backend/http/cookies"
	cookiePackage "github.com/tobigiwa/golang-security-backend/http/cookies"
	"github.com/tobigiwa/golang-security-backend/internal/store"
	"github.com/tobigiwa/golang-security-backend/pkg/logging"
	"golang.org/x/crypto/bcrypt"
)

// Webpp is application struct
type WebApp struct {
	Store  *store.Store
	Logger *logging.Logger
}

func (a *WebApp) generateHashedPassword(password string) ([]byte, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 15)
	return hashedPassword, err
}

func (a *WebApp) CheckRouteMethod(w http.ResponseWriter, r *http.Request, httpAllowRoute string) {
	if r.Method != httpAllowRoute {
		w.Header().Set("Allow", httpAllowRoute)
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

}

func (a *WebApp) CheckCookie(w http.ResponseWriter, r *http.Request, key string) string {
	secretKey, err := hex.DecodeString("JO/g8r/73iIJ6L8D7mBGU3pxKe5PMNNo3PS91hTWZRY=")
	if err != nil {
		log.Fatal(err)
	}
	value, err := cookies.ReadEncryptedCookie(r, key, secretKey)
	if err != nil {
		switch {
		case errors.Is(err, http.ErrNoCookie):
			a.ClientError(w, http.StatusBadRequest, "cookie not found")
		case errors.Is(err, cookiePackage.ErrInvalidValue):
			a.ClientError(w, http.StatusBadRequest, "invalid cookie")
		default:
			a.Logger.LogError(err, "APP")
			a.ServerError(w, "cookie invalid")
		}
	}
	return value
}

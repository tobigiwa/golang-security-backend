package http

import (
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/tobigiwa/golang-security-backend/http/cookies"
	cookiePackage "github.com/tobigiwa/golang-security-backend/http/cookies"
)

func (a *WebApp) AuthorizationBackend(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if !a.CheckRouteMethod(w, r, []string{http.MethodGet, http.MethodPost}) {
			return
		}

		cookie, _ := r.Cookie("llll")
		cookie.Valid()
		value := a.CheckCookie(w, r, "llll")
		user := a.DeserializeUserModel(value)
		fmt.Printf("\nUSER IS %T\n", user)
		fmt.Print(user.Email, "----", user.Username, "----", user.Status, "\n\n")
		http.SetCookie(w, cookie)
		next.ServeHTTP(w, r)

	})
}

func (a *WebApp) CheckCookie(w http.ResponseWriter, r *http.Request, key string) string {
	secretKey, err := hex.DecodeString("13d6b4dff8f84a10851021ec8608f814570d562c92fe6b5ec4c9f595bcb3234b")
	if err != nil {
		a.Logger.LogError(err, "APP")
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

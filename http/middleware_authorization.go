package http

import (
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"

	"github.com/tobigiwa/golang-security-backend/http/cookies"
	cookiePackage "github.com/tobigiwa/golang-security-backend/http/cookies"
)

func (a *WebApp) AuthorizationBackend(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		cookie, err := r.Cookie("cookie")
		if errors.Is(err, http.ErrNoCookie) {
			a.ClientError(w, http.StatusBadRequest, "cookie not found")
			return
		}

		value, err := a.retieveUserFromCookie(w, r, "cookie")
		if err != nil {
			a.Logger.LogError(err, "APP")
			return
		}
		user := a.DeserializeUserModel(value)

		fmt.Printf("\nUSER IS %T\n", user)
		fmt.Print(user.Email, "----", user.Username, "----", user.Status, "\n\n")

		http.SetCookie(w, cookie)
		next.ServeHTTP(w, r)

	})
}

func (a *WebApp) retieveUserFromCookie(w http.ResponseWriter, r *http.Request, key string) (string, error) {
	secretKey, err := hex.DecodeString("13d6b4dff8f84a10851021ec8608f814570d562c92fe6b5ec4c9f595bcb3234b")
	if err != nil {
		a.Logger.LogError(err, "APP")
	}
	value, err := cookies.ReadEncryptedCookie(r, secretKey)
	if err != nil {
		switch {
		case errors.Is(err, http.ErrNoCookie):
			a.ClientError(w, http.StatusBadRequest, "cookie not found")
			return "", err
		case errors.Is(err, cookiePackage.ErrInvalidValue):
			a.ClientError(w, http.StatusBadRequest, "invalid cookie")
			return "", err
		default:
			a.ServerError(w, "cookie invalid")
			return "", err
		}
	}
	return value, nil
}

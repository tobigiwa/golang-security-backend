package http

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	cookiePackage "github.com/tobigiwa/golang-security-backend/http/cookies"
	"golang.org/x/crypto/bcrypt"
)

func (a *WebApp) authenticationBackend(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		a.CheckRouteMethod(w, r, http.MethodPost)
		err := r.ParseForm()
		if err != nil {
			a.ClientError(w, http.StatusBadRequest, "invalid form data")
			return
		}
		email, password := r.PostForm.Get("email"), r.PostForm.Get("password")
		if email == "" || password == "" {
			a.ClientError(w, http.StatusBadRequest, "incomplete form data")
			return
		}
		status, err := a.Authenticate(email, password)
		if err != nil {
			if errors.Is(err, errInvalidCredentials) {
				a.ClientError(w, http.StatusForbidden, http.StatusText(http.StatusForbidden))
				return
			} else {
				a.ServerError(w, err.Error())
				return
			}
		}

		a.CreateCookie(w, status)
		fmt.Print()
		next.ServeHTTP(w, r)
	})
}

func (a *WebApp) Authenticate(email, password string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	hashedPassword, status, err := a.DbModel.FetchUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, errInvalidCredentials) {
			return "", errInvalidCredentials
		} else {
			return "", err
		}
	}
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return "", errInvalidCredentials
		} else {
			return "", err
		}
	}
	return status, nil

}

func (a *WebApp) CreateCookie(w http.ResponseWriter, payload string) {
	secretKey, err := hex.DecodeString("13d6b4dff8f84a10851021ec8608f814570d562c92fe6b5ec4c9f595bcb3234b")
	if err != nil {
		log.Fatal(err)
	}

	cookie := http.Cookie{
		Name:    "llll",
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
		return
	}
}

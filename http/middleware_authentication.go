package http

import (
	"encoding/hex"
	"errors"
	"log"
	"net/http"
	"time"

	cookiePackage "github.com/tobigiwa/golang-security-backend/http/cookies"
	"github.com/tobigiwa/golang-security-backend/internal/store"
)

func (a *WebApp) authenticationBackend(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		email, password := r.PostForm.Get("email"), r.PostForm.Get("password")
		user, err := a.GetUser(email, password)
		if err != nil {
			if errors.Is(err, store.ErrNotFound) {
				a.ClientError(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
				return
			}
		}
		err = a.Store.ValidateUserCredentials(user, password)
		if err != nil {
			a.ClientError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
			return
		}
		response := a.SerializeUserModel(&user)

		err = a.CreateCookie(w, response.String())
		if err != nil {
			a.ServerError(w, err.Error())
			a.Logger.LogError(err, "APP")
		}

		next.ServeHTTP(w, r)
	})
}

func (a *WebApp) GetUser(email, password string) (store.UserModel, error) {
	user, err := a.Store.FetchUser(email)
	if err != nil {
		return user, err
	}
	return user, nil

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

// if email == "" || password == "" {
// 	a.ClientError(w, http.StatusBadRequest, "incomplete form data")
// 	return
// }

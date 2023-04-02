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
		user, err := a.Authenticate(email, password)
		if err != nil {
			if errors.Is(err, errInvalidCredentials) {
				a.ClientError(w, http.StatusForbidden, http.StatusText(http.StatusForbidden))
				return
			} else {
				a.ServerError(w, err.Error())
				return
			}
		}

		response := a.SerializeUserModel(&user)
		a.CreateCookie(w, response.String())
		next.ServeHTTP(w, r)
	})
}

func (a *WebApp) Authenticate(email, password string) (store.UserModel, error) {
	user, err := a.Store.FetchUser(email)
	if err != nil {
		if errors.Is(err, errInvalidCredentials) {
			return store.UserModel{}, errInvalidCredentials
		} else {
			return store.UserModel{}, err
		}
	}
	// err = bcrypt.CompareHashAndPassword(user.Password, []byte(password))
	// if err != nil {
	// 	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
	// 		return store.UserModel{}, errInvalidCredentials
	// 	} else {
	// 		return store.UserModel{}, err
	// 	}
	// }
	return user, nil

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

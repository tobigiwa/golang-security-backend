package main

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/tobigiwa/golang-security-backend/internal/models"

	"golang.org/x/crypto/bcrypt"
)

func (a *WebApp) authenticate(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.Header().Set("Allow", http.MethodPost)
			w.Header().Add("Content-Type", "application/json")
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}

		err := r.ParseForm()
		if err != nil {
			a.clientError(w, http.StatusBadRequest)
			return
		}
		email, password := r.PostForm.Get("email"), r.PostForm.Get("password")
		err = a.authenticate(email, password)



		

func (a *WebApp) isAuthenticated(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		_, ok := r.Context().Value(isAuthenticatedContextKey).(bool)
		if !ok {
			a.clientError(w, http.StatusUnauthorized)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}


// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 		defer cancel()
// 		hashedPassword, err := a.dbModel.FetchUserByEmail(ctx, email)
// 		if err != nil {
// 			if errors.Is(err, models.ErrInvalidCredentials) {
// 				return models.ErrInvalidCredentials
// 			} else {
// 				return err
// 			}
// 		}
// 		err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
// 		if err != nil {
// 			if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
// 				return models.ErrInvalidCredentials
// 			} else {
// 				return err
// 			}
// 		}

// 		next.ServeHTTP(w, r)
// 	})
// }
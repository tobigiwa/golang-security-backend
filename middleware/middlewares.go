package middleware

import (
	"context"
	"errors"
	"net/http"

	"github.com/tobigiwa/golang-security-backend/internal/models"
)

type contextKey int

const (
	isAuthenticatedContextKey contextKey = iota + 1
)

type Middlewares func(http.HandlerFunc) http.HandlerFunc

func authenticationBackend(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.Header().Set("Allow", http.MethodPost)
			w.Header().Add("Content-Type", "application/json")
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}
		err := r.ParseForm()
		if err != nil {
			a.clientError(w, http.StatusBadRequest, "invalid form data")
			return
		}
		email, password := r.PostForm.Get("email"), r.PostForm.Get("password")
		if email == "" || password == "" {
			a.clientError(w, http.StatusBadRequest, "incomplete form data")
			return
		}
		err = a.authenticate(email, password)
		if err != nil {
			if errors.Is(err, models.ErrInvalidCredentials) {
				a.clientError(w, http.StatusForbidden, http.StatusText(http.StatusForbidden))
				return
			} else {
				a.undefinedError(w, err.Error())
				return
			}
		}
		ctx := context.WithValue(r.Context(), isAuthenticatedContextKey, true)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func (a *WebApp) isAuthenticated(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		_, ok := r.Context().Value(isAuthenticatedContextKey).(bool)
		if !ok {
			a.clientError(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
			// http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}

package http

import (
	"context"
	"errors"
	"net/http"
)

type contextKey int

const (
	isAuthenticatedContextKey contextKey = iota + 1
)

type Application interface {
	Authenticate(email, password string) error
	ClientError(w http.ResponseWriter, httpStatus int, text string)
	ServerError(w http.ResponseWriter, text string)
}

func authenticationBackend(a Application) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			if r.Method != http.MethodPost {
				w.Header().Set("Allow", http.MethodPost)
				w.Header().Add("Content-Type", "application/json")
				http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
				return
			}

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

			err = a.Authenticate(email, password)

			if err != nil {
				if errors.Is(err, errInvalidCredentials) {
					a.ClientError(w, http.StatusForbidden, http.StatusText(http.StatusForbidden))
					return
				} else {
					a.ServerError(w, err.Error())
					return
				}
			}

			ctx := context.WithValue(r.Context(), isAuthenticatedContextKey, true)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})

	}
}

// func isAuthenticated(next http.HandlerFunc) http.HandlerFunc {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

// 		_, ok := r.Context().Value(isAuthenticatedContextKey).(bool)
// 		if !ok {
// 			a.ClientError(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
// 			// http.Redirect(w, r, "/login", http.StatusSeeOther)
// 			return
// 		}
// 		next.ServeHTTP(w, r)
// 	})
// }

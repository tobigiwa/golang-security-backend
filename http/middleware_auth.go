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

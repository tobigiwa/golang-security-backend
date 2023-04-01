package http

import (
	"errors"
	"net/http"
)

func (a *WebApp) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "Close")
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				a.Logger.LogFatal(errors.New("got a panic"), "APP")
			}
		}()

		next.ServeHTTP(w, r)
	})
}

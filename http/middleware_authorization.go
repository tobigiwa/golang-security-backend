package http

import "net/http"

func (a *WebApp) AuthorizationBackend(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		a.CheckRouteMethod(w, r, http.MethodGet)
		
	})
}
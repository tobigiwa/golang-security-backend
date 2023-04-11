package http

import "net/http"

// routes returns defined routes on the muxtiplexer
func (a *WebApp) Routes() http.Handler {

	Login := http.HandlerFunc(a.Login)
	Welcome := http.HandlerFunc(a.Welcome)
	createUser := http.HandlerFunc(a.CreateUser)

	mux := http.NewServeMux()

	mux.HandleFunc("/", a.Home)
	mux.Handle("/welcome", a.httpMethod(http.MethodGet, a.AuthorizationBackend(Welcome)))
	mux.Handle("/login", a.httpMethod(http.MethodPost, a.authenticationBackend(Login)))
	mux.Handle("/createuser", a.httpMethod(http.MethodPost, createUser))

	return a.recoverPanic(a.logRequest(mux))

}

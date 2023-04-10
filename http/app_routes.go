package http

import "net/http"

// routes returns defined routes on the muxtiplexer
func (a *WebApp) Routes() http.Handler {

	Login := http.HandlerFunc(a.Login)
	Welcome := http.HandlerFunc(a.Welcome)
	Signup := http.HandlerFunc(a.Signup)

	mux := http.NewServeMux()

	mux.HandleFunc("/", a.Home)
	mux.Handle("/signup", a.httpMethod(http.MethodPost, Signup))
	mux.Handle("/login", a.httpMethod(http.MethodPost, a.authenticationBackend(Login)))
	mux.Handle("/welcome", a.httpMethod(http.MethodGet, a.AuthorizationBackend(Welcome)))

	return a.recoverPanic(a.logRequest(mux))

}

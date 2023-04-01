package http

import "net/http"

// routes returns defined routes on the muxtiplexer
func (a *WebApp) Routes() http.Handler {

	Login := http.HandlerFunc(a.Login)
	Welcome := http.HandlerFunc(a.Welcome)

	mux := http.NewServeMux()

	mux.HandleFunc("/", a.Home)
	mux.HandleFunc("/signup", a.Signup)
	mux.Handle("/login", a.authenticationBackend(Login))
	mux.Handle("/welcome", a.AuthorizationBackend(Welcome))

	return a.recoverPanic(mux)

}

package http

import "net/http"

// routes returns defined routes on the muxtiplexer
func (a *WebApp) Routes() *http.ServeMux {

	Login := http.HandlerFunc(a.Login)

	mux := http.NewServeMux()

	mux.HandleFunc("/", a.Home)
	mux.HandleFunc("/signup", a.Signup)
	mux.Handle("/login", a.authenticationBackend(Login))

	return mux

}

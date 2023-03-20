package main

import "net/http"

// routes returns defined routes on the muxtiplexer
func (a *WebApp) Routes() *http.ServeMux {

	mux := http.NewServeMux()
	mux.HandleFunc("/", a.Home)
	mux.HandleFunc("/signup", a.Signup)

	return mux

}

package main

import "net/http"

func (a *WebApp) routes() http.Handler {

	mux := http.NewServeMux()

	mux.HandleFunc("/signup", a.signup)

	return mux

}

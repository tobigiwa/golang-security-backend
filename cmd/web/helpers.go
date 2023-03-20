package main

import "net/http"

func (a *WebApp) clientError(w http.ResponseWriter, httpStatus int) {
	w.WriteHeader(httpStatus)
	http.Error(w, http.StatusText(httpStatus), httpStatus)
}

func (a *WebApp) serverError(w http.ResponseWriter, httpStatus int) {
	w.WriteHeader(httpStatus)
	http.Error(w, http.StatusText(httpStatus), httpStatus)
}

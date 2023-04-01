package http

import (
	"errors"
	"net/http"
)

func (a *WebApp) Home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome to the HomePage, allowed to everyone"))
}

func (a *WebApp) Signup(w http.ResponseWriter, r *http.Request) {
	if !a.CheckRouteMethod(w, r, []string{http.MethodPost}) {
		return
	}
	err := r.ParseForm()
	if err != nil {
		a.ClientError(w, http.StatusBadRequest, "invalid form data")
		return
	}
	email, username, password := r.PostForm.Get("email"), r.PostForm.Get("username"), r.PostForm.Get("password")
	if email == "" || password == "" || username == "" {
		a.ClientError(w, http.StatusBadRequest, "incomplete form data")
		return
	}
	hashedPassword, err := a.generateHashedPassword(password)
	if err != nil {
		a.ClientError(w, http.StatusUnsupportedMediaType, "password is of incorrect type: "+err.Error())
		return
	}
	err = a.Store.Insert(email, username, string(hashedPassword))
	if err != nil {
		if errors.Is(err, errDuplicateEmail) {
			a.ClientError(w, http.StatusConflict, "Email already used")
			return
		} else if errors.Is(err, errDuplicateUsername) {
			a.ClientError(w, http.StatusConflict, "Username already used")
			return
		} else {
			a.Logger.LogError(err, "APP")
			a.ServerError(w, err.Error())
			return
		}
	}
	w.Write([]byte("Signup SUCCESSFUL"))
	// http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (a *WebApp) Login(w http.ResponseWriter, r *http.Request) {

	w.Write([]byte("Login SUCCESSFUL"))
	// http.Redirect(w, r, "/welcome", http.StatusSeeOther)
}

func (a *WebApp) Welcome(w http.ResponseWriter, r *http.Request) {

	w.Write([]byte("WELCOME TO AUTHORIZED PAGE"))
}

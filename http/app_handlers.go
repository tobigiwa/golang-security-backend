package http

import (
	"errors"
	"net/http"

	"github.com/tobigiwa/golang-security-backend/internal/store"
)

func (a *WebApp) Home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome to the HomePage, allowed to everyone"))
}

func (a *WebApp) Signup(w http.ResponseWriter, r *http.Request) {

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
	err = a.Store.CreateUser(email, username, password)
	if err != nil {
		if errors.Is(err, store.ErrDuplicateEmail) {
			a.ClientError(w, http.StatusConflict, "Email already used")
			return
		} else if errors.Is(err, store.ErrDuplicateUsername) {
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

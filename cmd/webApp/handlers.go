package main

import (
	"errors"
	"net/http"
)

func (a *WebApp) Home(w http.ResponseWriter, r *http.Request) {

	w.Write([]byte("Welcome to the HomePage, allowed to everyone"))
}

func (a *WebApp) Signup(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		w.Header().Add("Content-Type", "application/json")
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		a.clientError(w, http.StatusBadRequest, "invalid form data")
		return
	}

	email, username, password := r.PostForm.Get("email"), r.PostForm.Get("username"), r.PostForm.Get("password")
	if email == "" || password == "" || username == "" {
		a.clientError(w, http.StatusBadRequest, "incomplete form data")
		return
	}

	hashedPassword, err := a.generateHashedPassword(password)
	if err != nil {
		a.clientError(w, http.StatusUnsupportedMediaType, "password is of incorrect type")
		w.Write([]byte(err.Error()))
		return
	}

	err = a.dbModel.Insert(email, username, string(hashedPassword))
	if err != nil {
		if errors.Is(err, errDuplicateEmail) {
			a.clientError(w, http.StatusConflict, "Email already used")
			return
		} else if errors.Is(err, errDuplicateUsername) {
			a.clientError(w, http.StatusConflict, "Username already used")
			return
		} else {
			a.serverError(w, err.Error())
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

	w.Write([]byte("WELCOME"))
}
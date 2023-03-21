package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/tobigiwa/golang-security-backend/internal/models"
)

type contextKey int

const (
	isAuthenticatedContextKey contextKey = iota + 1
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
		a.clientError(w, http.StatusBadRequest)
		return
	}

	email, username, password := r.PostForm.Get("email"), r.PostForm.Get("username"), r.PostForm.Get("password")
	hashedPassword, err := a.generateHashedPassword(password)
	if err != nil {
		a.clientError(w, http.StatusUnsupportedMediaType)
		w.Write([]byte(err.Error()))
		return
	}
	uuID := a.generateNewUUID()
	if uuID == "" {
		uuID = a.generateNewUUID()
	}

	err = a.dbModel.Insert(uuID, email, username, string(hashedPassword))
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			a.clientError(w, http.StatusConflict)
			fmt.Fprintln(w, "Email already already infmt.Fprintf(w, ) use")
			return
		} else if errors.Is(err, models.ErrDuplicateUsername) {
			a.clientError(w, http.StatusConflict)
			fmt.Fprintln(w, "Username already already in use")
			return
		} else {
			a.undefinedError(w, err.Error())
			return
		}
	}
	w.Write([]byte("INSERT WAS SUCCESSFUL"))

}

func (a *WebApp) Login(w http.ResponseWriter, r *http.Request) {
	
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			a.clientError(w, http.StatusForbidden)
			return
		} else {
			a.undefinedError(w, err.Error())
			return
		}
	}
	ctx := context.WithValue(r.Context(), isAuthenticatedContextKey, true)
	r = r.WithContext(ctx)
	http.Redirect(w, r, "/welcome", http.StatusSeeOther)
}

func (a *WebApp) Welcome(w http.ResponseWriter, r *http.Request) {

	w.Write([]byte("WELCOME"))
}

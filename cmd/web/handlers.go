package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/tobigiwa/golang-security-backend/internal/models"
)

type contextKey int

const (
	ContextKeyOne contextKey = iota + 1
)

func (a *WebApp) Home(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	fmt.Fprint(w, "Welcome to the HomePage")
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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = a.dbModel.Insert(ctx, email, username, string(hashedPassword))
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
	fmt.Fprintln(w, "INSERT WAS SUCCESSFUL")

}

func (a *WebApp) Login(w http.ResponseWriter, r *http.Request) {
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
	email, password := r.PostForm.Get("email"), r.PostForm.Get("password")
	id, err := a.authenticate(email, password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			a.clientError(w, http.StatusForbidden)
		} else {
			a.undefinedError(w, err.Error())
		}
	}
	ctx := context.WithValue(r.Context(), ContextKeyOne, id)
	r = r.WithContext(ctx)
	http.Redirect(w, r, "/welcome", http.StatusSeeOther)

}

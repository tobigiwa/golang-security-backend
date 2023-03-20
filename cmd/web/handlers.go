package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/tobigiwa/golang-security-backend/internal/models"
	"golang.org/x/crypto/bcrypt"
)

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
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 15)
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
			fmt.Fprintln(w, "Email already already in use")
			return
		} else if errors.Is(err, models.ErrDuplicateUsername) {
			fmt.Fprintln(w, "Username already already in use")
			return
		} else {
			fmt.Fprintln(w, err.Error())
			return
		}
	}
	fmt.Fprintln(w, "INSERT WAS SUCCESSFUL")

}

func (a *WebApp) Home(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	fmt.Fprint(w, "Welcome to the HomePage")
}

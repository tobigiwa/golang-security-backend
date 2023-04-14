package http

import (
	"errors"
	"net/http"

	"github.com/tobigiwa/golang-security-backend/internal/service"
)

func (a *WebApp) Home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome to the HomePage, allowed to everyone"))
}

func (a *WebApp) CreateUser(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		a.ClientError(w, http.StatusBadRequest, "invalid form data")
		return
	}
	email, username, password := r.PostForm.Get("email"), r.PostForm.Get("username"), r.PostForm.Get("password")

	err = a.Service.CreateUser(email, username, password)
	if err != nil {
		if errors.Is(err, service.ErrDuplicateEmail) {
			a.ClientError(w, http.StatusConflict, "Email already used")
			return
		} else if errors.Is(err, service.ErrDuplicateUsername) {
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

func (a *WebApp) CreateSUperUser(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		a.ClientError(w, http.StatusBadRequest, "invalid form data")
		return
	}
	email, username, password := r.PostForm.Get("email"), r.PostForm.Get("username"), r.PostForm.Get("password")

	err = a.Service.CreateSuperUser(email, username, password)
	if err != nil {
		if errors.Is(err, service.ErrDuplicateEmail) {
			a.ClientError(w, http.StatusConflict, "Email already used")
			return
		} else if errors.Is(err, service.ErrDuplicateUsername) {
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
	
	user, err := a.GetUser(r)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrNotFound):
			a.ClientError(w, http.StatusBadRequest, "user not found")
			return
		case errors.Is(err, service.ErrInvalidCredentials):
			a.ClientError(w, http.StatusBadRequest, "invalid user credentials")
			return
		default:
			a.ClientError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
			return
		}
	}
	serilizedUser := a.SerializeUserModel(&user)
	err = a.CreateCookie(w, serilizedUser.String())
	if err != nil {
		a.ServerError(w, err.Error())
		a.Logger.LogError(err, "APP")
	}

	w.Write([]byte("Login SUCCESSFUL"))

	// http.Redirect(w, r, "/welcome", http.StatusSeeOther)
}

func (a *WebApp) Welcome(w http.ResponseWriter, r *http.Request) {

	w.Write([]byte("WELCOME TO AUTHORIZED PAGE"))
}

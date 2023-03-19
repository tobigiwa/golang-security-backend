package main

import (
	"fmt"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func (a *WebApp) signup(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		fmt.Println(err)
	}
	email, username, password := r.PostForm.Get("email"), r.PostForm.Get("username"), r.PostForm.Get("password")
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 15)
	if err != nil {
		fmt.Println(err)
	}
	a.dbModel.Insert(email, username, string(hashedPassword))

}

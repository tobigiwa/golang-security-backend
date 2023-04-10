package http

import (
	"fmt"
	"net/http"
	"unicode/utf8"

	"github.com/go-playground/validator/v10"
)

func (a *WebApp) formValidator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		err := r.ParseForm()
		if err != nil {
			a.ClientError(w, http.StatusBadRequest, "invalid form data")
			return
		}
		email, username, password := r.PostForm.Get("email"), r.PostForm.Get("username"), r.PostForm.Get("password")
		u := UserValidator{
			Email:    email,
			Username: username,
			Password: password,
		}
		Validate(u)
		next.ServeHTTP(w, r)
	})
}

type UserValidator struct {
	Email    string `json:"email" validate:"email"`
	Username string `json:"name" validate:"username"`
	Password string `json:"password" validate:"passwd"`
}

func newValidator() *validator.Validate {
	return validator.New()
}

func ValidatePassword(fl validator.FieldLevel) bool {
	paswd := fl.Field().String()
	return utf8.RuneCountInString(paswd) > 7
}

func ValidateUsername(fl validator.FieldLevel) bool {
	username := fl.Field().String()
	return username != ""
}

func Validate(data UserValidator) {
	v := newValidator()
	v.RegisterValidation("passwd", ValidatePassword)
	v.RegisterValidation("username", ValidateUsername)

	err := v.Struct(data)
	if err != nil {
		validationErr, ok := err.(validator.ValidationErrors)
		if !ok {
			fmt.Println(err)
			return
		}
		for _, vErr := range validationErr {
			fmt.Printf("'%s' has a value of '%v' which does not satisfy '%s'.\n", vErr.Field(), vErr.Value(), vErr.Tag())
			switch vErr.Field() {
			case "Email":
				fmt.Println("Email must be a valid email format e.g example@example.com")
			case "Username":
				fmt.Println("Username must not be empty")
			case "Password":
				fmt.Println("Password must be 8 character or more")
			}
		}
	}
}

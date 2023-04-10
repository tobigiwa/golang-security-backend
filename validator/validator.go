package validator

import (
	"fmt"
	"unicode/utf8"

	"github.com/go-playground/validator/v10"
)

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

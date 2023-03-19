package main

type UserSignupForm struct {
	Email    string `form:"email"`
	Username string `form:"username"`
	Password string `from:"password"`
}

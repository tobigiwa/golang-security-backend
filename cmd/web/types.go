package main

type UserSignupForm struct {
	Email    string `form:"email"`
	Username string `form:"username"`
	Password string `from:"password"`
}

type Users struct {
	Email           string `json:"email"`
	Username        string `json:"username"`
	Hashed_password string `json:"passwor"`
	Status          string `json:"status"`
}

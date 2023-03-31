package http

import (
	"bytes"
	"encoding/gob"
	"strings"

	"github.com/tobigiwa/golang-security-backend/internal/store"
)

type UserResponseModel struct {
	Email    string
	Username string
	Status   string
}

func (a *WebApp) SerializeUserModel(user *store.UserModel) bytes.Buffer {
	var buf bytes.Buffer
	response := UserResponseModel{
		Email:    user.Email,
		Username: user.Username,
		Status:   user.Status,
	}
	err := gob.NewEncoder(&buf).Encode(&response)
	if err != nil {
		a.Logger.LogError(err, "DB")
	}
	return buf
}

func (a *WebApp) DeserializeUserModel(str string) store.UserModel {
	var user store.UserModel

	reader := strings.NewReader(str)
	err := gob.NewDecoder(reader).Decode(&user)
	if err != nil {
		a.Logger.LogError(err, "DB")
	}
	return user
}
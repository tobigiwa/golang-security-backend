package http

import (
	"context"
	"net/http"
	"time"

	"github.com/tobigiwa/golang-security-backend/internal/service"
)

func (a *WebApp) GetUser(r *http.Request) (service.UserModel, error) {

	email, password := r.PostForm.Get("email"), r.PostForm.Get("password")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	user, err := a.Service.GetAndValidateUser(ctx, email, password)
	if err != nil {
		return user, err
	}
	return user, nil
}

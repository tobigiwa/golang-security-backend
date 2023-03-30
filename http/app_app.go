package http

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/tobigiwa/golang-security-backend/internal/store"
	"github.com/tobigiwa/golang-security-backend/pkg/logging"
	"golang.org/x/crypto/bcrypt"
)

// Webpp is application struct
type WebApp struct {
	DbModel *store.UserModel
	Logger  *logging.Logger
}

func (a *WebApp) generateHashedPassword(password string) ([]byte, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 15)
	return hashedPassword, err
}

func (a *WebApp) Authenticate(email, password string) error {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	hashedPassword, err := a.DbModel.FetchUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, errInvalidCredentials) {
			return errInvalidCredentials
		} else {
			return err
		}
	}
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return errInvalidCredentials
		} else {
			return err
		}
	}
	return nil

}

func (a *WebApp) CheckRouteMethod(w http.ResponseWriter, r *http.Request, httpAllowRoute string) {
	if r.Method != httpAllowRoute {
		w.Header().Set("Allow", httpAllowRoute)
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

}



func (a *WebApp) securityCookie(r *http.Request) http.Cookie {
	cookie := http.Cookie{
		Name: "ddjjdjd",
		Value: "fdjdjd",
		Expires: time.Now().Add(30 * time.Minute),
		// MaxAge: ,
		Secure: true,
		HttpOnly: true,
		SameSite: ,


	}


}

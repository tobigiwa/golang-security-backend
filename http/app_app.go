package http

import (
	"net/http"

	"github.com/tobigiwa/golang-security-backend/internal/store"
	"github.com/tobigiwa/golang-security-backend/pkg/logging"
	"golang.org/x/crypto/bcrypt"
)

// Webpp is application struct
type WebApp struct {
	Store  *store.Store
	Logger *logging.Logger
}

func (a *WebApp) generateHashedPassword(password string) ([]byte, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 15)
	return hashedPassword, err
}

// CheckRouteMethod returns True if request.Method is NOT equal to httpAllowRoute and False otherwise.
// func (a *WebApp) CheckRouteMethod(w http.ResponseWriter, r *http.Request, httpAllowedRoutes []string) bool {
// 	seen_required := false
// 	for _, route := range httpAllowedRoutes {
// 		w.Header().Add("Allow", route)
// 		if !seen_required {
// 			seen_required = r.Method == route
// 		}
// 	}
// 	if seen_required {
// 		return true
// 	}

// 	http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
// 	return false

// }

func (a *WebApp) CheckRouteMethod(w http.ResponseWriter, r *http.Request, httpAllowedRoutes []string) bool {
	flag := false
	for _, route := range httpAllowedRoutes {
		w.Header().Add("Allow", route)
		if r.Method == route {
			flag = true
		}
	}
	if flag {
		return true
	}
	http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	return false
}

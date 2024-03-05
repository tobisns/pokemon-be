package auth

import (
	"errors"
	"learngo/pkg/services/users"
	"log"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/julienschmidt/httprouter"
)

type authenticator struct {
	Secret string
}

func New(secret string) *authenticator {
	return &authenticator{secret}
}

func (a *authenticator) Authorize(n httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		claims := &users.Claims{}
		token, err := r.Cookie("token")
		if err != nil {
			switch {
			case errors.Is(err, http.ErrNoCookie):
				http.Error(w, "cookie not found", http.StatusBadRequest)
			default:
				log.Println(err)
				http.Error(w, "server error", http.StatusInternalServerError)
			}
			return
		}

		jwtoken, err := jwt.ParseWithClaims(token.Value, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(a.Secret), nil
		})

		if err != nil {
			log.Println(err)
			http.Error(w, "server error.", http.StatusInternalServerError)
			return
		}

		if !jwtoken.Valid {
			http.Error(w, "invalid token.", http.StatusBadRequest)
			return
		}

		if !claims.IsAdmin {
			http.Error(w, "unauthorized request.", http.StatusBadRequest)
		}
		// call registered handler
		n(w, r, ps)
	}
}

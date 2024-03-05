package transport

import (
	"database/sql"
	"encoding/json"
	"errors"
	"learngo/pkg/services/users"
	"learngo/pkg/services/users/store"
	"learngo/pkg/utils/middlewares/cors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/julienschmidt/httprouter"
)

type JsonErrorResponse struct {
	Error *ApiError `json:"error"`
}

type ApiError struct {
	Status int16  `json:"status"`
	Title  string `json:"title"`
}

type handler struct {
	UserService users.Service
	Secret      string
}

// Activate sets all the services required for articles and registers all the endpoints with the engine.
func Activate(router *httprouter.Router, db *sql.DB, secret *string) {
	userService := users.New(store.New(db))
	newHandler(router, userService, secret)
}

func newHandler(router *httprouter.Router, as users.Service, secret *string) {
	h := handler{
		UserService: as,
		Secret:      *secret,
	}

	router.POST("/users", cors.MiddleCORS(h.Create))
	router.POST("/users/login", cors.MiddleCORS(h.Login))
	router.POST("/users/logout", cors.MiddleCORS(h.Logout))
	router.GET("/users/authenticate", cors.MiddleCORS(h.Authenticate))
}

func (h *handler) Create(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	var uc users.User
	var ur users.UserResponse

	if err := json.NewDecoder(r.Body).Decode(&uc); err != nil {
		status, message := handleError(err)
		http.Error(w, message.Error(), status)
		return
	}

	log.Printf("creating user %v", uc)
	ur, err := h.UserService.Create(r.Context(), &uc)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err != nil {
		status, message := handleError(err)
		http.Error(w, message.Error(), status)
		return
	}

	response := ur
	if err := json.NewEncoder(w).Encode(response); err != nil {
		panic(err)
	}
}

func (h *handler) Login(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	var uc users.User

	if err := json.NewDecoder(r.Body).Decode(&uc); err != nil {
		status, message := handleError(err)
		http.Error(w, message.Error(), status)
		return
	}

	log.Printf("authenticating user %v", uc)
	ut, err := h.UserService.LogIn(r.Context(), &uc, h.Secret)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if err != nil {
		status, message := handleError(err)
		http.Error(w, message.Error(), status)
		return
	}

	cookie := http.Cookie{
		Name:     "token",
		Value:    ut.Token,
		Path:     "/",
		MaxAge:   int(time.Until(ut.ExpirationTime).Seconds()),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(w, &cookie)
	response := ut.ExpirationTime
	if err := json.NewEncoder(w).Encode(response); err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
}

func (h *handler) Logout(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	// immediately clear the token cookie
	cookie := http.Cookie{
		Name:    "token",
		Value:   "",
		Path:    "/",
		Expires: time.Unix(0, 0),

		HttpOnly: true,
	}

	http.SetCookie(w, &cookie)
}

func (h *handler) Authenticate(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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
		return []byte(h.Secret), nil
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

	type Wrapper struct {
		IsAdmin bool `json:"is_admin"`
	}

	response := Wrapper{IsAdmin: claims.IsAdmin}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		panic(err)
	}
}

func handleError(e error) (int, error) {
	switch e {
	case users.ErrCreate:
		return http.StatusBadRequest, e
	case users.ErrAuth:
		return http.StatusBadRequest, e
	default:
		if numErr, ok := e.(*strconv.NumError); ok && numErr.Err == strconv.ErrSyntax {
			return http.StatusBadRequest, errors.New("bad request")
		}
		// Check if the error implements json.UnmarshalerError interface
		if _, ok := e.(*json.UnmarshalTypeError); ok {
			return http.StatusBadRequest, errors.New("bad request")
		}
		return http.StatusInternalServerError, e
	}
}

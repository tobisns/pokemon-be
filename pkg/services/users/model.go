package users

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type User struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserResponse struct {
	Username string `json:"username" binding:"required"`
	Message  string `json:"messages"`
}

type Claims struct {
	Username string `json:"username"`
	IsAdmin  bool   `json:"is_admin"`
	jwt.RegisteredClaims
}

type UserToken struct {
	Token          string    `josn:"token" binding:"required"`
	ExpirationTime time.Time `json:"expiration_time" binding:"required"`
}

type UserLoginResponse struct {
	IsAdmin        bool      `json:"is_admin"`
	ExpirationTime time.Time `json:"expiration_time"`
}

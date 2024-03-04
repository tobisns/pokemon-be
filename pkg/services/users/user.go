package users

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Repo defines the DB level interaction of pokemons
type Repo interface {
	Create(ctx context.Context, username, password string) (UserResponse, error)
	Authenticate(ctx context.Context, username, password string) (bool, error)
}

// Service defines the service level contract that other services
// outside this package can use to interact with Pokemon resources
type Service interface {
	Create(ctx context.Context, user *User) (UserResponse, error)
	LogIn(ctx context.Context, user *User, secret string) (UserToken, error)
}

type user struct {
	repo Repo
}

// New Service instance
func New(repo Repo) Service {
	return &user{repo}
}

func (s *user) Create(ctx context.Context, user *User) (UserResponse, error) {
	h := sha256.New()
	h.Write([]byte(user.Password))
	bs := h.Sum(nil)

	hashedPassword := hex.EncodeToString(bs)

	return s.repo.Create(ctx, user.Username, hashedPassword)
}

func (s *user) LogIn(ctx context.Context, user *User, secret string) (UserToken, error) {
	h := sha256.New()
	h.Write([]byte(user.Password))
	bs := h.Sum(nil)

	hashedPassword := hex.EncodeToString(bs)

	isAdmin, err := s.repo.Authenticate(ctx, user.Username, hashedPassword)
	if err != nil {
		return UserToken{}, err
	}

	expirationTime := time.Now().Add(10 * time.Minute)
	claims := &Claims{
		Username: user.Username,
		IsAdmin:  isAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return UserToken{}, err
	}

	return UserToken{Token: tokenString, ExpirationTime: expirationTime}, nil
}

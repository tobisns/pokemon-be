package users

import (
	"errors"
)

var (
	// ErrCreate ...
	ErrCreate = errors.New("data could not be created")
	ErrAuth   = errors.New("invalid credentials")
)

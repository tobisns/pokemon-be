package pokemons

import (
	"errors"
)

var (
	// ErrNotFound ...
	ErrNotFound = errors.New("requested data could not be found")

	// ErrQuery ...
	ErrQuery = errors.New("bad request")

	// ErrCreate ...
	ErrCreate = errors.New("data could not be created")

	// ErrUpdate ...
	ErrUpdate = errors.New("data could not be updated")

	// ErrPrepareDB ...
	ErrDB = errors.New("server error")
)

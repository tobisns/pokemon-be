package pokemons

import (
	"context"
)

// Repo defines the DB level interaction of pokemons
type Repo interface {
	Get(ctx context.Context, id int) (Pokemon, error)
	GetAll(ctx context.Context, limit, offset int) ([]Pokemon, error)
	Find(ctx context.Context, query string, limit, offset int) ([]Pokemon, error)
	Create(ctx context.Context, pr PokemonCreateUpdate) (int, error)
	Update(ctx context.Context, pr PokemonCreateUpdate, id int) error
}

// Service defines the service level contract that other services
// outside this package can use to interact with Pokemon resources
type Service interface {
	Get(ctx context.Context, id int) (Pokemon, error)
	GetAll(ctx context.Context, limit, offset int) ([]Pokemon, error)
	Find(ctx context.Context, query string, limit, offset int) ([]Pokemon, error)
	Create(ctx context.Context, pr PokemonCreateUpdate) (Pokemon, error)
	Update(ctx context.Context, pr PokemonCreateUpdate, id int) (Pokemon, error)
}

type pokemon struct {
	repo Repo
}

// New Service instance
func New(repo Repo) Service {
	return &pokemon{repo}
}

func (s *pokemon) Get(ctx context.Context, id int) (Pokemon, error) {
	return s.repo.Get(ctx, id)
}

func (s *pokemon) GetAll(ctx context.Context, limit, offset int) ([]Pokemon, error) {
	return s.repo.GetAll(ctx, limit, offset)
}

func (s *pokemon) Find(ctx context.Context, query string, limit, offset int) ([]Pokemon, error) {
	return s.repo.Find(ctx, query, limit, offset)
}

func (s *pokemon) Create(ctx context.Context, pr PokemonCreateUpdate) (Pokemon, error) {
	id, err := s.repo.Create(ctx, pr)
	if err != nil {
		return Pokemon{}, err
	}
	return s.repo.Get(ctx, id)
}

func (s *pokemon) Update(ctx context.Context, ar PokemonCreateUpdate, id int) (Pokemon, error) {
	if err := s.repo.Update(ctx, ar, id); err != nil {
		return Pokemon{}, err
	}
	return s.Get(ctx, id)
}

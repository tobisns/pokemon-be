package pokemons

import (
	"context"
)

// Repo defines the DB level interaction of pokemons
type Repo interface {
	Get(ctx context.Context, name string) (Pokemon, error)
	GetAll(ctx context.Context, limit, offset int) (Pokemons, error)
	Find(ctx context.Context, query string, limit, offset int) (Pokemons, error)
	Create(ctx context.Context, pr *PokemonCreateUpdate) (string, error)
	Update(ctx context.Context, name string, pr *PokemonCreateUpdate) error
	Delete(ctx context.Context, name string) (string, error)
	GetEvoTree(ctx context.Context, id int) (EvolutionTree, error)
	CreateEvoTree(ctx context.Context, ei *EvolutionCreate) (int, error)
	InsertToEvoTree(ctx context.Context, id int, ei *EvolutionCreate) (int, error)
	DeleteFromTree(ctx context.Context, id int, names *DeleteEvoData) (int, error)
	GetTypes(ctx context.Context) (Types, error)
	GetSameTypes(ctx context.Context, id int) (TypePokemonResponse, error)
}

// Service defines the service level contract that other services
// outside this package can use to interact with Pokemon resources
type Service interface {
	Get(ctx context.Context, name string) (Pokemon, error)
	GetAll(ctx context.Context, limit, offset int, query string) (Pokemons, error)
	Create(ctx context.Context, pr *PokemonCreateUpdate) (Pokemon, error)
	Update(ctx context.Context, name string, pr *PokemonCreateUpdate) (Pokemon, error)
	Delete(ctx context.Context, name string) (PokemonDeleteResponse, error)
	GetEvoTree(ctx context.Context, id int) (EvolutionTree, error)
	CreateEvoTree(ctx context.Context, ei *EvolutionCreate) (EvolutionTree, error)
	InsertToEvoTree(ctx context.Context, id int, ei *EvolutionCreate) (EvolutionTree, error)
	DeleteFromTree(ctx context.Context, id int, names *DeleteEvoData) (EvolutionTree, error)
	GetTypes(ctx context.Context) (Types, error)
	GetSameTypes(ctx context.Context, id int) (TypePokemonResponse, error)
}

type pokemon struct {
	repo Repo
}

// New Service instance
func New(repo Repo) Service {
	return &pokemon{repo}
}

func (s *pokemon) Get(ctx context.Context, name string) (Pokemon, error) {
	return s.repo.Get(ctx, name)
}

func (s *pokemon) GetAll(ctx context.Context, limit, offset int, query string) (Pokemons, error) {
	if query != "" {
		return s.repo.Find(ctx, query, limit, offset)
	}
	return s.repo.GetAll(ctx, limit, offset)
}

func (s *pokemon) Create(ctx context.Context, pr *PokemonCreateUpdate) (Pokemon, error) {
	name, err := s.repo.Create(ctx, pr)
	if err != nil {
		return Pokemon{}, err
	}
	return s.repo.Get(ctx, name)
}

func (s *pokemon) Update(ctx context.Context, name string, ar *PokemonCreateUpdate) (Pokemon, error) {
	if err := s.repo.Update(ctx, name, ar); err != nil {
		return Pokemon{}, err
	}
	return s.Get(ctx, ar.Name)
}

func (s *pokemon) Delete(ctx context.Context, name string) (PokemonDeleteResponse, error) {
	name, err := s.repo.Delete(ctx, name)
	if err != nil {
		return PokemonDeleteResponse{Name: name, Message: "delete failed."}, err
	}
	return PokemonDeleteResponse{Name: name, Message: "delete success."}, err
}

func (s *pokemon) GetEvoTree(ctx context.Context, id int) (EvolutionTree, error) {
	return s.repo.GetEvoTree(ctx, id)
}

func (s *pokemon) CreateEvoTree(ctx context.Context, ei *EvolutionCreate) (EvolutionTree, error) {
	id, err := s.repo.CreateEvoTree(ctx, ei)
	if err != nil {
		return EvolutionTree{}, err
	}
	return s.repo.GetEvoTree(ctx, id)
}

func (s *pokemon) InsertToEvoTree(ctx context.Context, id int, ei *EvolutionCreate) (EvolutionTree, error) {
	id, err := s.repo.InsertToEvoTree(ctx, id, ei)
	if err != nil {
		return EvolutionTree{}, err
	}
	return s.repo.GetEvoTree(ctx, id)
}

func (s *pokemon) DeleteFromTree(ctx context.Context, id int, names *DeleteEvoData) (EvolutionTree, error) {
	id, err := s.repo.DeleteFromTree(ctx, id, names)
	if err != nil {
		return EvolutionTree{}, err
	}

	et, err := s.repo.GetEvoTree(ctx, id)
	if err == ErrNotFound {
		return EvolutionTree{}, nil
	} else if err != nil {
		return EvolutionTree{}, err
	}

	return et, nil
}

func (s *pokemon) GetTypes(ctx context.Context) (Types, error) {
	return s.repo.GetTypes(ctx)
}

func (s *pokemon) GetSameTypes(ctx context.Context, id int) (TypePokemonResponse, error) {
	return s.repo.GetSameTypes(ctx, id)
}

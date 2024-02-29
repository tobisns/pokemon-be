package store

import (
	"context"
	"database/sql"
	"errors"
	"learngo/pkg/services/pokemons"
	"log"
)

const (
	selectPokemon      = `SELECT name, COALESCE(image_url, '') AS image_url, COALESCE(evo_tree_id, -1) AS evo_tree_id FROM pokemons WHERE name=$1`
	selectManyPokemons = `SELECT name, COALESCE(image_url, '') AS image_url, COALESCE(evo_tree_id, -1) AS evo_tree_id FROM pokemons LIMIT $1 OFFSET $2`
	searchPokemon      = `SELECT id, name, COALESCE(image_url, '') AS image_url, COALESCE(evo_tree_id, -1) AS evo_tree_id FROM pokemons WHERE name LIKE %$1% LIMIT $2 OFFSET $3`
	insertPokemon      = `INSERT INTO pokemons (name, image_url, evo_tree_id)`
	updatePokemonName  = `UPDATE pokemons SET name=$1, image_url=$2 WHERE id=$3`
	deletePokemon      = `DELETE FROM pokemons WHERE id=$1`
)

type pokemonRepo struct {
	DB *sql.DB
}

// New creates an instance of the accountRepo.
func New(conn *sql.DB) pokemons.Repo {
	return &pokemonRepo{conn}
}

// Get retrieves the article with the given id
func (r *pokemonRepo) Get(ctx context.Context, name string) (pokemons.Pokemon, error) {
	var pr pokemons.Pokemon

	err := r.DB.QueryRow(selectPokemon, name).
		Scan(&pr.Name, &pr.ImageUrl, &pr.EvolutionTree)
	if err != nil {
		log.Println(err)
		return pr, errors.New("error")
	}

	return pr, nil
}

func (r *pokemonRepo) GetAll(ctx context.Context, limit, offset int) ([]pokemons.Pokemon, error) {
	pl := make([]pokemons.Pokemon, 0)

	rows, err := r.DB.Query(selectManyPokemons, limit, offset)
	if err != nil {
		log.Println(ctx, "unable to query db: %s", err.Error())
		return pl, errors.New("error")
	}
	defer rows.Close()

	for rows.Next() {
		var pr pokemons.Pokemon
		if err := rows.Scan(&pr.Name, &pr.ImageUrl, &pr.EvolutionTree); err != nil {
			log.Println(ctx, "unable to scan db rows: %s", err.Error())
			return pl, errors.New("error")
		}

		pl = append(pl, pr)
	}

	return pl, nil
}
func (r *pokemonRepo) Find(ctx context.Context, query string, limit, offset int) ([]pokemons.Pokemon, error) {
	return nil, nil
}
func (r *pokemonRepo) Create(ctx context.Context, pr pokemons.PokemonCreateUpdate) (string, error) {
	return "", nil
}
func (r *pokemonRepo) Update(ctx context.Context, pr pokemons.PokemonCreateUpdate, name string) error {
	return nil
}

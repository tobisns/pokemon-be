package store

import (
	"context"
	"database/sql"
	"learngo/pkg/services/pokemons"
	"log"
)

const (
	selectPokemon       = `SELECT name, COALESCE(image_url, '') AS image_url, COALESCE(evo_tree_id, -1) AS evo_tree_id, COALESCE(height, 0) AS height, COALESCE(weight, 0) AS weight, COALESCE(hp, 0) as hp, COALESCE(atk, 0) AS atk, COALESCE(def, 0) AS def, COALESCE(sa, 0) AS sa, COALESCE(sd, 0) AS sd, COALESCE(spd, 0) AS spd FROM pokemons WHERE name=$1`
	selectManyPokemons  = `SELECT name, COALESCE(image_url, '') AS image_url FROM pokemons LIMIT $1 OFFSET $2`
	searchPokemon       = `SELECT name, COALESCE(image_url, '') AS image_url FROM pokemons WHERE name LIKE '%' || $1 || '%' LIMIT $2 OFFSET $3`
	insertPokemon       = `INSERT INTO pokemons (name, image_url, height, weight, hp, atk, def, sa, sd, spd) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING name`
	updatePokemon       = `UPDATE pokemons SET name = $1, image_url = $2, height = $3, weight = $4, hp = $5, atk = $6, def = $7, sa = $8, sd = $9, spd = $10 WHERE name = $11`
	deletePokemon       = `DELETE FROM pokemons WHERE name=$1`
	selectEvolutionTree = `SELECT id, level, pokemon_name FROM evo_tree WHERE id=$1`
	newTreeId           = `SELECT new_tree() AS id`
	insertEvolutionTree = `INSERT INTO evo_tree (id, pokemon_name, level) VALUES ($1, $2, $3)`
	deleteFromTree      = `DELETE FROM evo_tree WHERE id=$1 AND pokemon_name=$2`
	selectTypes         = `SELECT * FROM type`
	searchSameType      = `SELECT * FROM pokemon_type WHERE type_id=$1`
	insertType          = `INSERT INTO type (name) VALUES ($1) RETURNING name, id`
	insertPokemonType   = `INSERT INTO pokemon_type (pokemon, type_id) VALUES ($1, $2) RETURNING pokemon`
	selectPokemonTypes  = `SELECT pt.type_id, t.name FROM pokemon_type pt JOIN type t ON pt.type_id = t.id WHERE pt.pokemon=$1`
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
		Scan(&pr.Name, &pr.ImageUrl, &pr.EvolutionTreeID, &pr.Height, &pr.Weight, &pr.Stat.HealthPoint, &pr.Stat.Attack, &pr.Stat.Defense, &pr.Stat.SpecialAttack, &pr.Stat.SpecialDefense, &pr.Stat.Speed)
	if err != nil {
		log.Println(err)
		return pr, pokemons.ErrNotFound
	}

	return pr, nil
}

func (r *pokemonRepo) GetAll(ctx context.Context, limit, offset int) (pokemons.Pokemons, error) {
	var pl pokemons.Pokemons

	rows, err := r.DB.Query(selectManyPokemons, limit, offset)
	if err != nil {
		log.Println(ctx, "unable to query db: %s", err.Error())
		return pl, pokemons.ErrQuery
	}
	defer rows.Close()

	for rows.Next() {
		var pr pokemons.PokemonLight
		if err := rows.Scan(&pr.Name, &pr.ImageUrl); err != nil {
			log.Println(ctx, "unable to scan db rows: %s", err.Error())
			return pl, pokemons.ErrDB
		}

		pl.Pokemons = append(pl.Pokemons, pr)
	}

	return pl, nil
}

func (r *pokemonRepo) Find(ctx context.Context, query string, limit, offset int) (pokemons.Pokemons, error) {
	var pl pokemons.Pokemons

	rows, err := r.DB.Query(searchPokemon, query, limit, offset)
	if err != nil {
		log.Println(ctx, "unable to query db: %s", err.Error())
		return pl, pokemons.ErrQuery
	}
	defer rows.Close()

	for rows.Next() {
		var pr pokemons.PokemonLight
		if err := rows.Scan(&pr.Name, &pr.ImageUrl); err != nil {
			log.Println(ctx, "unable to scan db rows: %s", err.Error())
			return pl, pokemons.ErrDB
		}

		pl.Pokemons = append(pl.Pokemons, pr)
	}

	return pl, nil
}
func (r *pokemonRepo) Create(ctx context.Context, pr *pokemons.PokemonCreateUpdate) (string, error) {
	var name string
	if err := r.DB.QueryRow(insertPokemon, (*pr).Name, (*pr).ImageUrl, (*pr).Height, (*pr).Weight, (*pr).Stat.HealthPoint, (*pr).Stat.Attack, (*pr).Stat.Defense, (*pr).Stat.SpecialAttack, (*pr).Stat.SpecialDefense, (*pr).Stat.Speed).Scan(&name); err != nil {
		log.Println(ctx, "unable to create pokemon: %s", err.Error())
		return "", pokemons.ErrCreate
	}

	log.Println(ctx, "created pokemon with name=%s", name)
	return name, nil
}
func (r *pokemonRepo) Update(ctx context.Context, name string, pr *pokemons.PokemonCreateUpdate) error {
	_, err := r.DB.Exec(updatePokemon, (*pr).Name, (*pr).ImageUrl, (*pr).Height, (*pr).Weight, (*pr).Stat.HealthPoint, (*pr).Stat.Attack, (*pr).Stat.Defense, (*pr).Stat.SpecialAttack, (*pr).Stat.SpecialDefense, (*pr).Stat.Speed, name)
	if err != nil {
		log.Println(ctx, "unable to update pokemon (%s): %s", pr.Name, err.Error())
		return pokemons.ErrUpdate
	}
	return nil
}

func (r *pokemonRepo) Delete(ctx context.Context, name string) (string, error) {
	result, err := r.DB.Exec(deletePokemon, name)
	if err != nil {
		log.Println(ctx, "unable to delete pokemon: %s", err.Error())
		return "", pokemons.ErrDB
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Println(ctx, "error getting rows affected: %s", err.Error())
		return "", pokemons.ErrDB
	}

	if rowsAffected == 0 {
		log.Println(ctx, "no pokemon with name=%s found to delete", name)
		return "", pokemons.ErrNotFound
	}

	log.Println(ctx, "deleted pokemon with name=%s", name)
	return name, nil
}

func (r *pokemonRepo) GetEvoTree(ctx context.Context, id int) (pokemons.EvolutionTree, error) {
	var et pokemons.EvolutionTree

	rows, err := r.DB.Query(selectEvolutionTree, id)
	if err != nil {
		log.Println(ctx, "unable to query db: %s", err.Error())
		return et, pokemons.ErrQuery
	}
	defer rows.Close()

	for rows.Next() {
		var er pokemons.EvolutionData
		if err := rows.Scan(&er.ID, &er.Level, &er.PokemonName); err != nil {
			log.Println(ctx, "unable to scan db rows: %s", err.Error())
			return et, pokemons.ErrCreate
		}

		et.EvolutionData = append(et.EvolutionData, er)
	}

	if et.EvolutionData == nil {
		return et, pokemons.ErrNotFound
	}

	return et, nil
}

func (r *pokemonRepo) CreateEvoTree(ctx context.Context, ei *pokemons.EvolutionCreate) (int, error) {
	var treeId int

	tx, err := r.DB.Begin()
	if err != nil {
		log.Println(ctx, "unable to begin transaction: %s", err.Error())
		return -1, pokemons.ErrDB
	}

	tx.QueryRow(newTreeId).Scan(&treeId)
	for _, data := range ei.EvolutionCreateData {
		_, err := tx.Exec(insertEvolutionTree, treeId, data.PokemonName, data.Level)
		if err != nil {
			tx.Rollback()
			log.Println(ctx, "unable to complete transaction: %s", err.Error())
			return -1, pokemons.ErrCreate
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Println(ctx, "unable to complete transaction: %s", err.Error())
		return -1, pokemons.ErrDB
	}

	return treeId, nil
}

func (r *pokemonRepo) InsertToEvoTree(ctx context.Context, id int, ei *pokemons.EvolutionCreate) (int, error) {
	tx, err := r.DB.Begin()
	if err != nil {
		log.Println(ctx, "unable to begin transaction: %s", err.Error())
		return -1, pokemons.ErrDB
	}

	for _, data := range ei.EvolutionCreateData {
		_, err := tx.Exec(insertEvolutionTree, id, data.PokemonName, data.Level)
		if err != nil {
			tx.Rollback()
			log.Println(ctx, "unable to complete transaction: %s", err.Error())
			return -1, pokemons.ErrCreate
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Println(ctx, "unable to complete transaction: %s", err.Error())
		return -1, pokemons.ErrDB
	}

	return id, nil
}

func (r *pokemonRepo) DeleteFromTree(ctx context.Context, id int, names *pokemons.DeleteEvoData) (int, error) {
	tx, err := r.DB.Begin()
	if err != nil {
		log.Println(ctx, "unable to begin transaction: %s", err.Error())
		return -1, pokemons.ErrDB
	}

	for _, pokemon := range names.Pokemons {
		_, err := tx.Exec(deleteFromTree, id, pokemon.Name)
		if err != nil {
			tx.Rollback()
			log.Println(ctx, "unable to complete transaction: %s", err.Error())
			return -1, pokemons.ErrNotFound
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Println(ctx, "unable to complete transaction: %s", err.Error())
		return -1, pokemons.ErrDB
	}

	return id, nil
}

func (r *pokemonRepo) GetTypes(ctx context.Context) (pokemons.Types, error) {
	var tl pokemons.Types

	rows, err := r.DB.Query(selectTypes)
	if err != nil {
		log.Println(ctx, "unable to query db: %s", err.Error())
		return tl, pokemons.ErrQuery
	}
	defer rows.Close()

	for rows.Next() {
		var tr pokemons.Type
		if err := rows.Scan(&tr.ID, &tr.Name); err != nil {
			log.Println(ctx, "unable to scan db rows: %s", err.Error())
			return tl, pokemons.ErrDB
		}

		tl.Types = append(tl.Types, tr)
	}

	return tl, nil
}

func (r *pokemonRepo) GetSameTypes(ctx context.Context, id int) (pokemons.TypePokemonResponse, error) {
	var pl pokemons.TypePokemonResponse

	rows, err := r.DB.Query(searchSameType, id)
	if err != nil {
		log.Println(ctx, "unable to query db: %s", err.Error())
		return pl, pokemons.ErrQuery
	}
	defer rows.Close()

	for rows.Next() {
		var pr pokemons.TypePokemon
		if err := rows.Scan(&pr.Name, &pr.TypeId); err != nil {
			log.Println(ctx, "unable to scan db rows: %s", err.Error())
			return pl, pokemons.ErrCreate
		}

		pl.Pokemons = append(pl.Pokemons, pr)
	}

	return pl, nil
}

func (r *pokemonRepo) CreateType(ctx context.Context, name string) (pokemons.Type, error) {
	var pt pokemons.Type
	if err := r.DB.QueryRow(insertType, name).Scan(&pt.Name, &pt.ID); err != nil {
		log.Println(ctx, "unable to create type: %s", err.Error())
		return pokemons.Type{}, pokemons.ErrCreate
	}

	log.Println(ctx, "created type with name=%s", name)
	return pt, nil
}

func (r *pokemonRepo) AssignType(ctx context.Context, name string, typeId int) (string, error) {
	var pn string
	if err := r.DB.QueryRow(insertPokemonType, name, typeId).Scan(&pn); err != nil {
		log.Println(ctx, "unable to assign type: %s", err.Error())
		return "", pokemons.ErrCreate
	}

	log.Println(ctx, "assigned type to name=%s", pn)
	return pn, nil
}

func (r *pokemonRepo) GetPokemonTypes(ctx context.Context, name string) (pokemons.Types, error) {
	var tl pokemons.Types

	rows, err := r.DB.Query(selectPokemonTypes, name)
	if err != nil {
		log.Println(ctx, "unable to query db: %s", err.Error())
		return tl, pokemons.ErrQuery
	}
	defer rows.Close()

	for rows.Next() {
		var tr pokemons.Type
		if err := rows.Scan(&tr.ID, &tr.Name); err != nil {
			log.Println(ctx, "unable to scan db rows: %s", err.Error())
			return tl, pokemons.ErrDB
		}

		tl.Types = append(tl.Types, tr)
	}

	return tl, nil
}

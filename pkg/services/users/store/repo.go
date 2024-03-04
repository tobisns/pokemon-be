package store

import (
	"context"
	"database/sql"
	"learngo/pkg/services/users"
	"log"
)

const (
	createUser = `INSERT INTO users (username, password) VALUES ($1, $2) RETURNING username`
	selectUser = `SELECT admin FROM users WHERE username=$1 AND password=$2`
)

type userRepo struct {
	DB *sql.DB
}

func New(conn *sql.DB) users.Repo {
	return &userRepo{conn}
}

func (r *userRepo) Create(ctx context.Context, username, password string) (users.UserResponse, error) {
	var ur users.UserResponse
	err := r.DB.QueryRow(createUser, username, password).Scan(&ur.Username)
	if err != nil {
		log.Println(ctx, "unable to create user: %s", err.Error())
		return ur, users.ErrCreate
	}

	ur.Message = "account created."
	return ur, nil
}

func (r *userRepo) Authenticate(ctx context.Context, username, password string) (bool, error) {
	var isAdmin bool

	err := r.DB.QueryRow(selectUser, username, password).Scan(&isAdmin)
	if err != nil {
		log.Println(ctx, "unable to create user: %s", err.Error())
		return false, users.ErrAuth
	}

	return isAdmin, nil
}

package store

import (
	"context"
	"errors"
	"learngo/pkg/services/albums"
)

// album represents data about a record album.

type albumRepo struct {
	DB *[]albums.Album
}

func New(dat *[]albums.Album) albums.Repo {
	return &albumRepo{dat}
}

func (r *albumRepo) Get(ctx context.Context, id string) (albums.Album, error) {
	for _, album := range *r.DB {
		if album.ID == id {
			return album, nil
		}
	}
	return albums.Album{}, errors.New("error")
}

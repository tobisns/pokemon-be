package api

import (
	"learngo/pkg/services/albums"
	album "learngo/pkg/services/albums/transport"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func Start() {
	// albums slice to seed record album data.
	var albumsData = []albums.Album{
		{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
		{ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
		{ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
	}

	router := httprouter.New()

	album.Activate(router, &albumsData)

	log.Fatal(http.ListenAndServe(":8080", router))
}

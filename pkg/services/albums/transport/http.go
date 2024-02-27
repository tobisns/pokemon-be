package transport

import (
	"encoding/json"
	"learngo/pkg/services/albums"
	"learngo/pkg/services/albums/store"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type JsonErrorResponse struct {
	Error *ApiError `json:"error"`
}

type ApiError struct {
	Status int16  `json:"status"`
	Title  string `json:"title"`
}

type handler struct {
	AlbumService albums.Service
}

// Activate sets all the services required for articles and registers all the endpoints with the engine.
func Activate(router *httprouter.Router, db *[]albums.Album) {
	albumService := albums.New(store.New(db))
	newHandler(router, albumService)
}

func newHandler(router *httprouter.Router, as albums.Service) {
	h := handler{
		AlbumService: as,
	}

	router.GET("/albums/:id", h.Get)
}

func (h *handler) Get(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	article, err := h.AlbumService.Get(r.Context(), params.ByName("id"))
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err != nil {
		// No album with the id in the url has been found
		w.WriteHeader(http.StatusNotFound)
		response := JsonErrorResponse{Error: &ApiError{Status: 404, Title: "Record Not Found"}}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			panic(err)
		}
	}

	response := article
	if err := json.NewEncoder(w).Encode(response); err != nil {
		panic(err)
	}
}

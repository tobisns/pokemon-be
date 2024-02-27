package transport

import (
	"database/sql"
	"encoding/json"
	"learngo/pkg/services/pokemons"
	"learngo/pkg/services/pokemons/store"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/schema"
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
	PokemonService pokemons.Service
}

// Activate sets all the services required for articles and registers all the endpoints with the engine.
func Activate(router *httprouter.Router, db *sql.DB) {
	pokemonService := pokemons.New(store.New(db))
	newHandler(router, pokemonService)
}

func newHandler(router *httprouter.Router, as pokemons.Service) {
	h := handler{
		PokemonService: as,
	}

	router.GET("/pokemons/:id", h.Get)
	router.GET("/pokemons", h.GetAll)
}

func (h *handler) Get(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	id, err := strconv.Atoi(params.ByName("id"))
	pokemon, err := h.PokemonService.Get(r.Context(), id)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err != nil {
		// No album with the id in the url has been found
		w.WriteHeader(http.StatusNotFound)
		response := JsonErrorResponse{Error: &ApiError{Status: 404, Title: "Record Not Found"}}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			panic(err)
		}
	}

	response := pokemon
	if err := json.NewEncoder(w).Encode(response); err != nil {
		panic(err)
	}
}

func (h *handler) GetAll(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	q := struct {
		Limit  int `form:"limit"`
		Offset int `form:"offset"`
	}{Limit: 25, Offset: 0}

	ctx := r.Context()
	if err := schema.NewDecoder().Decode(&q, r.URL.Query()); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("retrieving all articles: offset=%d limit=%d", q.Limit, q.Offset)
	pokemons, err := h.PokemonService.GetAll(ctx, q.Limit, q.Offset)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		response := JsonErrorResponse{Error: &ApiError{Status: 404, Title: "Record Not Found"}}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			panic(err)
		}
	}

	response := pokemons
	if err := json.NewEncoder(w).Encode(response); err != nil {
		panic(err)
	}
}

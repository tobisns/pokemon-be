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

	router.GET("/pokemons/:name", h.Get)
	router.GET("/pokemons", h.GetAll)
	router.POST("/pokemons", h.Create)
	router.PUT("/pokemons/:name", h.Update)
	router.DELETE("/pokemons/:name", h.Delete)
	router.GET("/evolution_tree/:id", h.GetEvoTree)
	router.POST("/evolution_tree", h.CreateEvoTree)
	router.PUT("/evolution_tree/:id", h.InsertToEvoTree)
}

func (h *handler) Get(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	pokemon, err := h.PokemonService.Get(r.Context(), params.ByName("name"))
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
	}

	response := pokemon
	if err := json.NewEncoder(w).Encode(response); err != nil {
		panic(err)
	}
}

func (h *handler) GetAll(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	q := struct {
		Query  string `form:"query"`
		Limit  int    `form:"limit"`
		Offset int    `form:"offset"`
	}{Limit: 25, Offset: 0}

	ctx := r.Context()
	if err := schema.NewDecoder().Decode(&q, r.URL.Query()); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("retrieving all articles: offset=%d limit=%d query=%s", q.Offset, q.Limit, q.Query)
	pokemons, err := h.PokemonService.GetAll(ctx, q.Limit, q.Offset, q.Query)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
	}

	response := pokemons
	if err := json.NewEncoder(w).Encode(response); err != nil {
		panic(err)
	}
}

func (h *handler) Create(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	var pc pokemons.PokemonCreateUpdate

	if err := json.NewDecoder(r.Body).Decode(&pc); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("creating pokemon %v", pc)
	pokemon, err := h.PokemonService.Create(r.Context(), &pc)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
	}

	response := pokemon
	if err := json.NewEncoder(w).Encode(response); err != nil {
		panic(err)
	}
}

func (h *handler) Update(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	var pc pokemons.PokemonCreateUpdate

	if err := json.NewDecoder(r.Body).Decode(&pc); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("updating pokemon %v", pc)
	pokemon, err := h.PokemonService.Update(r.Context(), params.ByName("name"), &pc)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
	}

	response := pokemon
	if err := json.NewEncoder(w).Encode(response); err != nil {
		panic(err)
	}
}

func (h *handler) Delete(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	pokemon, err := h.PokemonService.Delete(r.Context(), params.ByName("name"))
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
	}

	response := pokemon
	if err := json.NewEncoder(w).Encode(response); err != nil {
		panic(err)
	}
}

func (h *handler) GetEvoTree(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	evoTree, err := h.PokemonService.GetEvoTree(r.Context(), id)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
	}

	response := evoTree
	if err := json.NewEncoder(w).Encode(response); err != nil {
		panic(err)
	}
}

func (h *handler) CreateEvoTree(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	var ec pokemons.EvolutionCreate

	if err := json.NewDecoder(r.Body).Decode(&ec); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("updating pokemon %v", ec)
	evoTree, err := h.PokemonService.CreateEvoTree(r.Context(), &ec)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
	}

	response := evoTree
	if err := json.NewEncoder(w).Encode(response); err != nil {
		panic(err)
	}
}

func (h *handler) InsertToEvoTree(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	var ec pokemons.EvolutionCreate
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	if err := json.NewDecoder(r.Body).Decode(&ec); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("updating pokemon %v", ec)
	evoTree, err := h.PokemonService.InsertToEvoTree(r.Context(), id, &ec)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
	}

	response := evoTree
	if err := json.NewEncoder(w).Encode(response); err != nil {
		panic(err)
	}
}

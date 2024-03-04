package transport

import (
	"database/sql"
	"encoding/json"
	"errors"
	"learngo/pkg/services/pokemons"
	"learngo/pkg/services/pokemons/store"
	auth "learngo/pkg/utils/middlewares"
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
	Secret         string
}

// Activate sets all the services required for articles and registers all the endpoints with the engine.
func Activate(router *httprouter.Router, db *sql.DB, secret *string) {
	pokemonService := pokemons.New(store.New(db))
	newHandler(router, pokemonService, secret)
}

func newHandler(router *httprouter.Router, as pokemons.Service, secret *string) {
	h := handler{
		PokemonService: as,
		Secret:         *secret,
	}

	auth := auth.New(h.Secret)

	router.GET("/pokemons/:name", h.Get)
	router.GET("/pokemons", h.GetAll)
	router.POST("/pokemons", auth.Authenticate(h.Create))
	router.PUT("/pokemons/:name", auth.Authenticate(h.Update))
	router.DELETE("/pokemons/:name", auth.Authenticate(h.Delete))
	router.GET("/evolution_tree/:id", h.GetEvoTree)
	router.POST("/evolution_tree", auth.Authenticate(h.CreateEvoTree))
	router.PUT("/evolution_tree/:id", auth.Authenticate(h.InsertToEvoTree))
	router.DELETE("/evolution_tree/:id", auth.Authenticate(h.DeleteFromEvoTree))
	router.GET("/types", h.GetTypes)
	router.GET("/types/:id", h.GetSameTypes)
}

func (h *handler) Get(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	pokemon, err := h.PokemonService.Get(r.Context(), params.ByName("name"))
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err != nil {
		status, message := handleError(err)
		http.Error(w, message.Error(), status)
		return
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
		status, message := handleError(err)
		http.Error(w, message.Error(), status)
		return
	}

	log.Printf("retrieving all articles: offset=%d limit=%d query=%s", q.Offset, q.Limit, q.Query)
	pokemons, err := h.PokemonService.GetAll(ctx, q.Limit, q.Offset, q.Query)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err != nil {
		status, message := handleError(err)
		http.Error(w, message.Error(), status)
		return
	}

	response := pokemons
	if err := json.NewEncoder(w).Encode(response); err != nil {
		panic(err)
	}
}

func (h *handler) Create(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	var pc pokemons.PokemonCreateUpdate

	if err := json.NewDecoder(r.Body).Decode(&pc); err != nil {
		status, message := handleError(err)
		http.Error(w, message.Error(), status)
		return
	}

	log.Printf("creating pokemon %v", pc)
	pokemon, err := h.PokemonService.Create(r.Context(), &pc)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err != nil {
		status, message := handleError(err)
		http.Error(w, message.Error(), status)
		return
	}

	response := pokemon
	if err := json.NewEncoder(w).Encode(response); err != nil {
		panic(err)
	}
}

func (h *handler) Update(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	var pc pokemons.PokemonCreateUpdate

	if err := json.NewDecoder(r.Body).Decode(&pc); err != nil {
		status, message := handleError(err)
		http.Error(w, message.Error(), status)
		return
	}

	log.Printf("updating pokemon %v", pc)
	pokemon, err := h.PokemonService.Update(r.Context(), params.ByName("name"), &pc)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err != nil {
		status, message := handleError(err)
		http.Error(w, message.Error(), status)
		return
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
		status, message := handleError(err)
		http.Error(w, message.Error(), status)
		return
	}

	response := pokemon
	if err := json.NewEncoder(w).Encode(response); err != nil {
		panic(err)
	}
}

func (h *handler) GetEvoTree(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		status, message := handleError(err)
		http.Error(w, message.Error(), status)
		return
	}
	evoTree, err := h.PokemonService.GetEvoTree(r.Context(), id)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err != nil {
		status, message := handleError(err)
		http.Error(w, message.Error(), status)
		return
	}

	response := evoTree
	if err := json.NewEncoder(w).Encode(response); err != nil {
		panic(err)
	}
}

func (h *handler) CreateEvoTree(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	var ec pokemons.EvolutionCreate

	if err := json.NewDecoder(r.Body).Decode(&ec); err != nil {
		status, message := handleError(err)
		http.Error(w, message.Error(), status)
		return
	}

	log.Printf("updating pokemon %v", ec)
	evoTree, err := h.PokemonService.CreateEvoTree(r.Context(), &ec)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err != nil {
		status, message := handleError(err)
		http.Error(w, message.Error(), status)
		return
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
		status, message := handleError(err)
		http.Error(w, message.Error(), status)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&ec); err != nil {
		status, message := handleError(err)
		http.Error(w, message.Error(), status)
		return
	}

	log.Printf("updating pokemon %v", ec)
	evoTree, err := h.PokemonService.InsertToEvoTree(r.Context(), id, &ec)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err != nil {
		status, message := handleError(err)
		http.Error(w, message.Error(), status)
		return
	}

	response := evoTree
	if err := json.NewEncoder(w).Encode(response); err != nil {
		panic(err)
	}
}

func (h *handler) DeleteFromEvoTree(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	var ed pokemons.DeleteEvoData
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		status, message := handleError(err)
		http.Error(w, message.Error(), status)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&ed); err != nil {
		status, message := handleError(err)
		http.Error(w, message.Error(), status)
		return
	}

	log.Printf("updating pokemon %v", ed)
	evoTree, err := h.PokemonService.DeleteFromTree(r.Context(), id, &ed)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err != nil {
		status, message := handleError(err)
		http.Error(w, message.Error(), status)
		return
	}

	response := evoTree
	if err := json.NewEncoder(w).Encode(response); err != nil {
		panic(err)
	}
}

func (h *handler) GetTypes(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	types, err := h.PokemonService.GetTypes(r.Context())
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err != nil {
		status, message := handleError(err)
		http.Error(w, message.Error(), status)
		return
	}

	response := types
	if err := json.NewEncoder(w).Encode(response); err != nil {
		panic(err)
	}
}

func (h *handler) GetSameTypes(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		status, message := handleError(err)
		http.Error(w, message.Error(), status)
		return
	}

	types, err := h.PokemonService.GetSameTypes(r.Context(), id)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err != nil {
		status, message := handleError(err)
		http.Error(w, message.Error(), status)
		return
	}

	response := types
	if err := json.NewEncoder(w).Encode(response); err != nil {
		panic(err)
	}
}

func handleError(e error) (int, error) {
	switch e {
	case pokemons.ErrNotFound:
		return http.StatusNotFound, e
	case pokemons.ErrUpdate:
		fallthrough
	case pokemons.ErrCreate:
		return http.StatusInternalServerError, e
	case pokemons.ErrQuery:
		return http.StatusBadRequest, e
	default:
		if numErr, ok := e.(*strconv.NumError); ok && numErr.Err == strconv.ErrSyntax {
			return http.StatusBadRequest, errors.New("bad request")
		}
		// Check if the error implements json.UnmarshalerError interface
		if _, ok := e.(*json.UnmarshalTypeError); ok {
			return http.StatusBadRequest, errors.New("bad request")
		}
		return http.StatusInternalServerError, e
	}
}

package api

import (
	"learngo/pkg/db"
	pokemon "learngo/pkg/services/pokemons/transport"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// Config defines what the API requires to run
type Config struct {
	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string
}

// Start initializes the API server, adding the reuired middleware and dependent services
func Start(cfg *Config) {
	conn, err := db.GetConnection(
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if conn != nil {
			conn.Close()
		}
	}()

	router := httprouter.New()

	pokemon.Activate(router, conn)

	log.Fatal(http.ListenAndServe(":8080", router))
}

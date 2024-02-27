package main

import (
	"learngo/pkg/api"
)

func main() {
	api.Start(&api.Config{
		DBHost:     "pokemon_db",
		DBPort:     5432,
		DBUser:     "postgres",
		DBPassword: "lololol",
		DBName:     "pokemon_db",
	})
}

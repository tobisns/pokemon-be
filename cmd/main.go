package main

import (
	"learngo/pkg/api"
)

func main() {
	api.Start(&api.Config{
		DBHost:     "localhost",
		DBPort:     5434,
		DBUser:     "postgres",
		DBPassword: "lololol",
		DBName:     "pokemon_db",
	})
}

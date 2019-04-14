package main

import (
	"log"

	"github.com/tac44/drift/environment"
	"github.com/tac44/drift/postgres"
)

func main() {
	environment.LoadEnv()

	database := postgres.NewDBConnection()
	defer database.CloseConnection()

	log.Println("Connected to drift!")

	exists := database.CheckForDriftMigrationsTable()

	if !exists {
		success := database.CreateDriftMigrationsTable()
		if !success {
			panic("Could not create drift migrations table")
		}
	}
}

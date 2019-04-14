package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/tac44/drift/postgres"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found...")
	}

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

package main

import (
	"log"

	"github.com/tac44/drift"
	"github.com/tac44/drift/environment"
	"github.com/tac44/drift/postgres"
)

func init() {
	environment.LoadEnv()
}

func main() {
	db, err := postgres.NewDBConnection()
	if err != nil {
		log.Fatal(err)
	}
	defer db.CloseConnection()

	log.Println("Connected to drift!")

	err = drift.Initialize(db)
	if err != nil {
		log.Fatal(err)
	}
}

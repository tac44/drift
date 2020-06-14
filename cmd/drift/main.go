package main

import (
	"github.com/tomcase/drift"
	"github.com/tomcase/drift/environment"
	"github.com/tomcase/drift/postgres"
	"log"
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

package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

// AppDB represents the db to apply migrations to
type AppDB struct {
	ctx context.Context
	db  *sql.DB
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found...")
	}

	host, _ := os.LookupEnv("DRIFT_PG_HOST")
	port, _ := os.LookupEnv("DRIFT_PG_PORT")
	user, _ := os.LookupEnv("DRIFT_PG_USER")
	pass, _ := os.LookupEnv("DRIFT_PG_PASS")
	dbName, _ := os.LookupEnv("DRIFT_PG_DB")

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, pass, dbName)

	log.Printf("psqlInfo: %s \n", psqlInfo)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	database := AppDB{
		ctx: context.Background(),
		db:  db,
	}

	log.Println("Connected to drift!")

	err = database.checkForDriftMigrationsTable()
	if err != nil && err != sql.ErrNoRows {
		panic(err)
	}

	if err == sql.ErrNoRows {
		err = database.createDriftMigrationsTable()
		if err != nil {
			log.Fatal(err)
		}
	}
}

func (db *AppDB) createDriftMigrationsTable() (err error) {

	log.Println("Creating drift_migrations table...")

	_, err = db.db.Exec(`
		CREATE TABLE "drift_migrations" (
		"id" uuid NOT NULL,
		"sequence" integer NOT NULL,
		"description" character varying(255) NOT NULL,
		"applied" date NOT NULL
		);
	`)

	if err != nil {
		log.Println("Created drift_migrations table!")
	}

	return err
}

func (db *AppDB) checkForDriftMigrationsTable() (err error) {
	query := `

		SELECT table_name
		FROM information_schema.tables
		WHERE table_catalog = 'drift' 
		AND table_name = 'drift_migrations'
		LIMIT 1;
	`

	log.Println(query)

	var tableName string
	err = db.db.QueryRow(query).Scan(&tableName)
	switch err {
	case sql.ErrNoRows:
		log.Println("drift_migrations table doesn't exist for database drift!")
	default:
		log.Println("drift_migrations table exists!")
	}

	return err
}

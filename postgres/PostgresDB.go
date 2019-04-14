package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	// Blank import necessary to use specify postgres drivers
	_ "github.com/lib/pq"
)

// DB represents the db to apply migrations to
type DB struct {
	ctx context.Context
	db  *sql.DB
}

// NewDBConnection returns a new DB connection
func NewDBConnection() *DB {
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

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	database := DB{
		ctx: context.Background(),
		db:  db,
	}

	return &database
}

// CloseConnection will close the connection to the database
func (db *DB) CloseConnection() {
	db.db.Close()
}

// CreateDriftMigrationsTable creates the migrations table
func (db *DB) CreateDriftMigrationsTable() (success bool) {

	log.Println("Creating drift_migrations table...")

	exec := `
		CREATE TABLE "drift_migrations" (
		"id" uuid NOT NULL,
		"sequence" integer NOT NULL,
		"description" character varying(255) NOT NULL,
		"applied" date NOT NULL
		);
	`

	log.Println(exec)

	_, err := db.db.Exec(exec)

	if err == nil {
		success = true
		log.Println("Created drift_migrations table!")
	}

	return success
}

// CheckForDriftMigrationsTable checks for the drift migrations table
func (db *DB) CheckForDriftMigrationsTable() (exists bool) {
	query := `
		SELECT table_name
		FROM information_schema.tables
		WHERE table_catalog = 'drift' 
		AND table_name = 'drift_migrations'
		LIMIT 1;
	`

	log.Println(query)

	var tableName string
	row := db.db.QueryRow(query)

	err := row.Scan(&tableName)
	if err != nil && err != sql.ErrNoRows {
		log.Fatal(err)
	}

	if row != nil {
		exists = true
		if err == sql.ErrNoRows {
			exists = false
			log.Println("drift_migrations table doesn't exist for database drift!")
		}
	}

	return exists
}

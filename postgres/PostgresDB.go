package postgres

import (
	"context"
	"database/sql"
	"fmt"
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
func NewDBConnection() (*DB, error) {
	host, _ := os.LookupEnv("DRIFT_PG_HOST")
	port, _ := os.LookupEnv("DRIFT_PG_PORT")
	user, _ := os.LookupEnv("DRIFT_PG_USER")
	pass, _ := os.LookupEnv("DRIFT_PG_PASS")
	dbName, _ := os.LookupEnv("DRIFT_PG_DB")

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, pass, dbName)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	database := DB{
		ctx: context.Background(),
		db:  db,
	}

	return &database, nil
}

// CloseConnection will close the connection to the database
func (db *DB) CloseConnection() error {
	err := db.db.Close()
	if err != nil {
		return err
	}
	return nil
}

// CreateDriftMigrationsTable creates the migrations table
func (db *DB) CreateDriftMigrationsTable() error {
	exec := `
		CREATE TABLE "drift_migrations" (
		"id" uuid NOT NULL,
		"sequence" integer NOT NULL,
		"description" character varying(255) NOT NULL,
		"applied" date NOT NULL
		);
	`
	_, err := db.db.Exec(exec)
	if err != nil {
		return err
	}

	return nil
}

// CheckForDriftMigrationsTable checks for the drift migrations table
func (db *DB) CheckForDriftMigrationsTable() (bool, error) {
	query := `
		SELECT table_name
		FROM information_schema.tables
		WHERE table_catalog = 'drift' 
		AND table_name = 'drift_migrations'
		LIMIT 1;
	`

	var tableName string
	row := db.db.QueryRow(query)

	err := row.Scan(&tableName)
	if err != nil && err != sql.ErrNoRows {
		return false, err
	}

	if row != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
	}

	return true, nil
}

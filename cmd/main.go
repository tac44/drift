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

	message, err := database.checkForDriftMigrationsTable()
	if err != nil {
		if err != sql.ErrNoRows {
			panic(err)
		}
	}

	log.Println(message)

	if err == sql.ErrNoRows {
		message, err = database.createDriftMigrationsTable()
		if err != nil {
			panic(err)
		}

		log.Println(message)
	}

}

func (db *AppDB) createDriftMigrationsTable() (message string, err error) {
	return "Creating drift_migrations table...", nil
}

func (db *AppDB) checkForDriftMigrationsTable() (message string, err error) {
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
	if err == nil {
		message = "drift_migrations table exists!"
	}
	if err == sql.ErrNoRows {
		message = "drift_migrations table doesn't exist for database drift!"
	}

	return message, err
}

package environment

import (
	"log"

	"github.com/joho/godotenv"
)

// LoadEnv will use godotenv to load environment variables
func LoadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found...")
	}
}

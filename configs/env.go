package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var err = godotenv.Load()

func EnvMongoURI() string {
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	return os.Getenv("MONGOURI")
}

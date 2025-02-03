package configs

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

var err = godotenv.Load()

func EnvMongoURI() string {
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	return os.Getenv("MONGOURI")
}

func GetSecretKey() string {
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	return os.Getenv("SECRET_KEY")
}

func GetExpirationTime() int {
	expirationDate, err := strconv.Atoi(os.Getenv("EXPIRATION_TIME"))
	if err != nil {
		return 60
	}
	return expirationDate
}

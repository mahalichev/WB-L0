package inits

import (
	"log"

	"github.com/joho/godotenv"
)

func LoadEnvironment() {
	if godotenv.Load() != nil {
		log.Fatal("failed to load .env")
	}
}

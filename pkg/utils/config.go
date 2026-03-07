package utils

import (
	"log"

	"github.com/joho/godotenv"
)

func LoadEnvironmentVariables() {
	if err := godotenv.Load(); err != nil {
		log.Printf(".env not loaded: %v", err)
	}
}

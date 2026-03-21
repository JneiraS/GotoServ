package utils

import (
	"errors"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnvironmentVariables() {
	if err := godotenv.Load(); err != nil {
		log.Printf(".env not loaded: %v", err)
	}
}

// ValidateSecretKey vérifie que SECRET_KEY est définie et a une longueur minimale de 16 caractères.
// Le programme s'arrête immédiatement si la clé est absente ou insuffisamment longue.
func ValidateSecretKey() {
	if err := checkSecretKey(os.Getenv("SECRET_KEY")); err != nil {
		log.Fatalf("SECRET_KEY invalide: %v", err)
	}
}

func checkSecretKey(key string) error {
	if key == "" {
		return errors.New("SECRET_KEY est vide, définissez-la dans .env ou via variable d'environnement")
	}
	if len(key) < 16 {
		return errors.New("SECRET_KEY trop courte (minimum 16 caractères)")
	}
	return nil
}

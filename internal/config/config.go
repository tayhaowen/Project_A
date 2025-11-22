package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port      string
	SecretJWT string
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		if os.IsNotExist(err) {
			log.Println(".env file not found, relying on environment variables")
		} else {
			return nil, fmt.Errorf("error loading .env file: %w", err)
		}
	} else {
		log.Println("Config loaded from .env")
	}

	return &Config{
		Port:      os.Getenv("PORT"),
		SecretJWT: os.Getenv("SECRET_JWT"),
	}, nil
}

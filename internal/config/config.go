package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Port      string
	SecretJWT string
	CacheTTL  time.Duration
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
		Port:      resolvePort(),
		SecretJWT: os.Getenv("SECRET_JWT"),
		CacheTTL:  resolveCacheTTL(),
	}, nil
}

const (
	defaultPort    = "8080"
	defaultCacheTTl = 30 * time.Minute
)

func resolvePort() string {
	if port := os.Getenv("PORT"); port != "" {
		return port
	}
	log.Printf("PORT not set, defaulting to %s", defaultPort)
	return defaultPort
}

func resolveCacheTTL() time.Duration {
	value := os.Getenv("CACHE_TTL_MINUTES")
	if value == "" {
		return defaultCacheTTl
	}
	minutes, err := strconv.Atoi(value)
	if err != nil || minutes <= 0 {
		log.Printf("Invalid CACHE_TTL_MINUTES=%q, using default %s", value, defaultCacheTTl)
		return defaultCacheTTl
	}
	return time.Duration(minutes) * time.Minute
}

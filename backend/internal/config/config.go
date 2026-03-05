package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL string
	Port        string
	Env         string
	JWTSecret   string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	return &Config{
		DatabaseURL: getEnv("DATABASE_URL", "postgres://capcurve:capcurve_dev@localhost:5432/capcurve_development?sslmode=disable"),
		Port:        getEnv("PORT", "8080"),
		Env:         getEnv("ENV", "development"),
		JWTSecret:   getEnv("JWT_SECRET", "dev_secret_change_in_production"),
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

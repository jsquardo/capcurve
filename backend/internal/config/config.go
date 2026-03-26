package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL  string
	Port         string
	Env          string
	JWTSecret    string
	AdminSecret  string
	SyncEnabled  bool
	SyncHour     int
	SyncMinute   int
	SyncWeekday  int
	SyncTimeZone string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	return &Config{
		DatabaseURL:  getEnv("DATABASE_URL", "postgres://capcurve:capcurve_dev@localhost:5432/capcurve_development?sslmode=disable"),
		Port:         getEnv("PORT", "8080"),
		Env:          getEnv("ENV", "development"),
		JWTSecret:    getEnv("JWT_SECRET", "dev_secret_change_in_production"),
		AdminSecret:  getEnv("ADMIN_SECRET", ""),
		SyncEnabled:  getEnvBool("SYNC_ENABLED", true),
		SyncHour:     getEnvInt("SYNC_HOUR", 5),
		SyncMinute:   getEnvInt("SYNC_MINUTE", 0),
		SyncWeekday:  getEnvInt("SYNC_WEEKDAY", int(1)),
		SyncTimeZone: getEnv("SYNC_TIMEZONE", "America/New_York"),
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func getEnvBool(key string, fallback bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return fallback
	}

	return parsed
}

func getEnvInt(key string, fallback int) int {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}

	return parsed
}

package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config holds all application configuration loaded from environment variables.
type Config struct {
	DBHost    string
	DBPort    string
	DBUser    string
	DBPass    string
	DBName    string
	JWTSecret string
	Port      string
}

// Load reads configuration from a .env file (if present) and then from
// environment variables, which always take precedence.
func Load() *Config {
	// Attempt to load .env; non-fatal if absent (env vars already set).
	if err := godotenv.Load(); err != nil {
		log.Println("[config] No .env file found, using environment variables")
	}

	cfg := &Config{
		DBHost:    getEnv("DB_HOST", "localhost"),
		DBPort:    getEnv("DB_PORT", "5432"),
		DBUser:    getEnv("DB_USER", "postgres"),
		DBPass:    getEnv("DB_PASS", "postgres"),
		DBName:    getEnv("DB_NAME", "coworking_spaces"),
		JWTSecret: getEnv("JWT_SECRET", "coworking_jwt_secret_2024_secure_key"),
		Port:      getEnv("PORT", "8002"),
	}

	return cfg
}

// getEnv returns the environment variable value or a fallback default.
func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

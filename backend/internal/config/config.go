package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

const (
	DEFAULT_PORT    = 8080
	DEFAULT_DB_PATH = "../data/app.db"
)

// Config holds all the configuration for the application
type Config struct {
	Port      int
	DBPath    string
	JWTSecret string
}

// LoadConfig loads the environment variables into the Config struct
func LoadConfig() *Config {
	// Load .env file if it exists.
	// We ignore the error because in production, variables might be set on the host system directly.
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("FATAL: JWT_SECRET environment variable is required but not set. Please provide a secure secret (e.g., run 'openssl rand -hex 32' to generate one).")
	}

	return &Config{
		Port:      getEnvAsInt("PORT", DEFAULT_PORT),
		DBPath:    getEnv("DB_PATH", DEFAULT_DB_PATH),
		JWTSecret: jwtSecret,
	}
}

// Helper function to read an environment variable or return a fallback string
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	fmt.Printf("Using default value for %s -> %s\n", key, fallback)
	return fallback
}

// Helper function to read an environment variable and convert it to an integer
func getEnvAsInt(key string, fallback int) int {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		fmt.Printf("Using default value for %s -> %d\n", key, fallback)
		return fallback
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		log.Printf("Warning: Invalid integer for %s, using fallback %d", key, fallback)
		return fallback
	}
	return value
}

package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfigDefaults(t *testing.T) {
	// Save and clear env vars
	savedPort := os.Getenv("PORT")
	savedDBPath := os.Getenv("DB_PATH")
	savedJWT := os.Getenv("JWT_SECRET")
	os.Unsetenv("PORT")
	os.Unsetenv("DB_PATH")
	os.Unsetenv("JWT_SECRET")
	defer func() {
		os.Setenv("PORT", savedPort)
		os.Setenv("DB_PATH", savedDBPath)
		os.Setenv("JWT_SECRET", savedJWT)
	}()

	cfg := LoadConfig()

	assert.Equal(t, DEFAULT_PORT, cfg.Port)
	assert.Equal(t, DEFAULT_DB_PATH, cfg.DBPath)
	assert.Equal(t, DEFAULT_JWT_SECRET, cfg.JWTSecret)
}

func TestLoadConfigWithEnv(t *testing.T) {
	savedPort := os.Getenv("PORT")
	savedDBPath := os.Getenv("DB_PATH")
	savedJWT := os.Getenv("JWT_SECRET")
	os.Setenv("PORT", "9090")
	os.Setenv("DB_PATH", "/custom/path/db.sqlite")
	os.Setenv("JWT_SECRET", "custom-secret")
	defer func() {
		os.Setenv("PORT", savedPort)
		os.Setenv("DB_PATH", savedDBPath)
		os.Setenv("JWT_SECRET", savedJWT)
	}()

	cfg := LoadConfig()

	assert.Equal(t, 9090, cfg.Port)
	assert.Equal(t, "/custom/path/db.sqlite", cfg.DBPath)
	assert.Equal(t, "custom-secret", cfg.JWTSecret)
}

func TestGetEnv(t *testing.T) {
	key := "TEST_GETENV_KEY"
	defer os.Unsetenv(key)

	result := getEnv(key, "fallback")
	assert.Equal(t, "fallback", result)

	os.Setenv(key, "custom")
	result = getEnv(key, "fallback")
	assert.Equal(t, "custom", result)
}

func TestGetEnvAsInt(t *testing.T) {
	key := "TEST_GETENVINT_KEY"
	defer os.Unsetenv(key)

	result := getEnvAsInt(key, 42)
	assert.Equal(t, 42, result)

	os.Setenv(key, "100")
	result = getEnvAsInt(key, 42)
	assert.Equal(t, 100, result)

	os.Setenv(key, "invalid")
	result = getEnvAsInt(key, 42)
	assert.Equal(t, 42, result)
}

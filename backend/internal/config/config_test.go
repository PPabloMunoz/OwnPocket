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
	_ = os.Unsetenv("PORT")
	_ = os.Unsetenv("DB_PATH")
	_ = os.Setenv("JWT_SECRET", "test-secret")
	defer func() {
		_ = os.Setenv("PORT", savedPort)
		_ = os.Setenv("DB_PATH", savedDBPath)
		_ = os.Setenv("JWT_SECRET", savedJWT)
	}()

	cfg := LoadConfig()

	assert.Equal(t, DEFAULT_PORT, cfg.Port)
	assert.Equal(t, DEFAULT_DB_PATH, cfg.DBPath)
	assert.Equal(t, "test-secret", cfg.JWTSecret)
}

func TestLoadConfigWithEnv(t *testing.T) {
	savedPort := os.Getenv("PORT")
	savedDBPath := os.Getenv("DB_PATH")
	savedJWT := os.Getenv("JWT_SECRET")
	_ = os.Setenv("PORT", "9090")
	_ = os.Setenv("DB_PATH", "/custom/path/db.sqlite")
	_ = os.Setenv("JWT_SECRET", "custom-secret")
	defer func() {
		_ = os.Setenv("PORT", savedPort)
		_ = os.Setenv("DB_PATH", savedDBPath)
		_ = os.Setenv("JWT_SECRET", savedJWT)
	}()

	cfg := LoadConfig()

	assert.Equal(t, 9090, cfg.Port)
	assert.Equal(t, "/custom/path/db.sqlite", cfg.DBPath)
	assert.Equal(t, "custom-secret", cfg.JWTSecret)
}

func TestGetEnv(t *testing.T) {
	key := "TEST_GETENV_KEY"
	defer func() { _ = os.Unsetenv(key) }()

	result := getEnv(key, "fallback")
	assert.Equal(t, "fallback", result)

	_ = os.Setenv(key, "custom")
	result = getEnv(key, "fallback")
	assert.Equal(t, "custom", result)
}

func TestGetEnvAsInt(t *testing.T) {
	key := "TEST_GETENVINT_KEY"
	defer func() { _ = os.Unsetenv(key) }()

	result := getEnvAsInt(key, 42)
	assert.Equal(t, 42, result)

	_ = os.Setenv(key, "100")
	result = getEnvAsInt(key, 42)
	assert.Equal(t, 100, result)

	_ = os.Setenv(key, "invalid")
	result = getEnvAsInt(key, 42)
	assert.Equal(t, 42, result)
}

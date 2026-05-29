package handler

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHealth(t *testing.T) {
	r, _, _ := setupTestRouter(t)
	w := executeRequest(r, "GET", "/api/v1/health", nil)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "ok", w.Body.String())
}

func TestRegister_Success(t *testing.T) {
	r, _, _ := setupTestRouter(t)
	w := executeRequest(r, "POST", "/api/v1/auth/register", map[string]any{
		"username": "newuser",
		"password": "password123",
		"email":    "new@example.com",
	})
	assert.Equal(t, http.StatusCreated, w.Code)

	resp := parseResponse(t, w)
	assert.Equal(t, "User registered successfully", resp["data"].(map[string]any)["message"])
	assert.Nil(t, resp["error"])
}

func TestRegister_Validation(t *testing.T) {
	r, _, _ := setupTestRouter(t)
	w := executeRequest(r, "POST", "/api/v1/auth/register", map[string]any{
		"username": "newuser",
		"password": "123",
	})
	assert.Equal(t, http.StatusBadRequest, w.Code)

	resp := parseResponse(t, w)
	assert.Nil(t, resp["data"])
	assert.NotEmpty(t, resp["error"])
}

func TestRegister_Duplicate(t *testing.T) {
	r, _, db := setupTestRouter(t)
	createTestUser(t, db)

	w := executeRequest(r, "POST", "/api/v1/auth/register", map[string]any{
		"username": "testuser",
		"password": "password123",
	})
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestLogin_Success(t *testing.T) {
	r, _, _ := setupTestRouter(t)

	// Register first
	executeRequest(r, "POST", "/api/v1/auth/register", map[string]any{
		"username": "testuser",
		"password": "password123",
		"email":    "test@example.com",
	})

	w := executeRequest(r, "POST", "/api/v1/auth/login", map[string]any{
		"username": "testuser",
		"password": "password123",
	})
	assert.Equal(t, http.StatusOK, w.Code)

	resp := parseResponse(t, w)
	assert.NotEmpty(t, resp["data"].(map[string]any)["token"])
}

func TestLogin_InvalidCredentials(t *testing.T) {
	r, _, _ := setupTestRouter(t)
	w := executeRequest(r, "POST", "/api/v1/auth/login", map[string]any{
		"username": "nonexistent",
		"password": "wrong",
	})
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	resp := parseResponse(t, w)
	assert.Equal(t, "invalid credentials", resp["error"])
}

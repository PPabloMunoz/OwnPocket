package handler

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateCategory_Success(t *testing.T) {
	r, _, db := setupTestRouter(t)
	createTestUser(t, db)

	w := executeRequest(r, "POST", "/api/v1/categories", map[string]any{
		"name": "Food",
		"type": "expense",
	})
	assert.Equal(t, http.StatusCreated, w.Code)

	resp := parseResponse(t, w)
	data := resp["data"].(map[string]any)
	assert.Equal(t, "Food", data["name"])
	assert.Equal(t, "expense", data["type"])
}

func TestCreateCategory_Validation(t *testing.T) {
	r, _, db := setupTestRouter(t)
	createTestUser(t, db)

	w := executeRequest(r, "POST", "/api/v1/categories", map[string]any{
		"name": "Bad",
		"type": "invalid",
	})
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetCategories(t *testing.T) {
	r, _, db := setupTestRouter(t)
	createTestUser(t, db)

	executeRequest(r, "POST", "/api/v1/categories", map[string]any{
		"name": "Food",
		"type": "expense",
	})
	executeRequest(r, "POST", "/api/v1/categories", map[string]any{
		"name": "Salary",
		"type": "income",
	})

	w := executeRequest(r, "GET", "/api/v1/categories", nil)
	assert.Equal(t, http.StatusOK, w.Code)

	resp := parseResponse(t, w)
	categories := resp["data"].([]any)
	assert.Len(t, categories, 2)
}

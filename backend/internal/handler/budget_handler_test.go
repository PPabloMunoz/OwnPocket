package handler

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateBudget_Success(t *testing.T) {
	r, _, db := setupTestRouter(t)
	createTestUser(t, db)

	// Create a category first
	catResp := executeRequest(r, "POST", "/api/v1/categories", map[string]any{
		"name": "Food",
		"type": "expense",
	})
	cat := parseResponse(t, catResp)
	catID := cat["data"].(map[string]any)["id"].(float64)

	w := executeRequest(r, "POST", "/api/v1/budgets", map[string]any{
		"category_id": catID,
		"period":      "2026-06",
		"amount":      500.00,
	})
	assert.Equal(t, http.StatusCreated, w.Code)

	resp := parseResponse(t, w)
	data := resp["data"].(map[string]any)
	assert.Equal(t, catID, data["category_id"].(float64))
	assert.Equal(t, "2026-06", data["period"])
}

func TestCreateBudget_Validation(t *testing.T) {
	r, _, db := setupTestRouter(t)
	createTestUser(t, db)

	w := executeRequest(r, "POST", "/api/v1/budgets", map[string]any{
		"category_id": 0,
		"period":      "",
		"amount":      0,
	})
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetBudgets(t *testing.T) {
	r, _, db := setupTestRouter(t)
	createTestUser(t, db)

	catResp := executeRequest(r, "POST", "/api/v1/categories", map[string]any{
		"name": "Food",
		"type": "expense",
	})
	cat := parseResponse(t, catResp)
	catID := cat["data"].(map[string]any)["id"].(float64)

	executeRequest(r, "POST", "/api/v1/budgets", map[string]any{
		"category_id": catID,
		"period":      "2026-06",
		"amount":      500.00,
	})

	w := executeRequest(r, "GET", "/api/v1/budgets", nil)
	assert.Equal(t, http.StatusOK, w.Code)

	resp := parseResponse(t, w)
	budgets := resp["data"].([]any)
	assert.Len(t, budgets, 1)
}

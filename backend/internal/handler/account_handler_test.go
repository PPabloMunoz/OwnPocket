package handler

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateAccount_Success(t *testing.T) {
	r, _, db := setupTestRouter(t)
	createTestUser(t, db)

	w := executeRequest(r, "POST", "/api/v1/accounts", map[string]any{
		"name": "Checking",
		"type": "checking",
	})
	assert.Equal(t, http.StatusCreated, w.Code)

	resp := parseResponse(t, w)
	data := resp["data"].(map[string]any)
	assert.Equal(t, "Checking", data["name"])
	assert.Equal(t, "checking", data["type"])
}

func TestCreateAccount_Validation(t *testing.T) {
	r, _, db := setupTestRouter(t)
	createTestUser(t, db)

	w := executeRequest(r, "POST", "/api/v1/accounts", map[string]any{
		"name": "Bad",
		"type": "invalid_type",
	})
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetAccounts(t *testing.T) {
	r, _, db := setupTestRouter(t)
	createTestUser(t, db)

	executeRequest(r, "POST", "/api/v1/accounts", map[string]any{
		"name": "Checking",
		"type": "checking",
	})
	executeRequest(r, "POST", "/api/v1/accounts", map[string]any{
		"name": "Savings",
		"type": "savings",
	})

	w := executeRequest(r, "GET", "/api/v1/accounts", nil)
	assert.Equal(t, http.StatusOK, w.Code)

	resp := parseResponse(t, w)
	accounts := resp["data"].([]any)
	assert.Len(t, accounts, 2)
}

func TestGetAccount_Success(t *testing.T) {
	r, _, db := setupTestRouter(t)
	createTestUser(t, db)

	createResp := executeRequest(r, "POST", "/api/v1/accounts", map[string]any{
		"name": "Checking",
		"type": "checking",
	})
	created := parseResponse(t, createResp)
	accountID := int(created["data"].(map[string]any)["id"].(float64))

	w := executeRequest(r, "GET", fmt.Sprintf("/api/v1/accounts/%d", accountID), nil)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetAccount_NotFound(t *testing.T) {
	r, _, db := setupTestRouter(t)
	createTestUser(t, db)

	w := executeRequest(r, "GET", "/api/v1/accounts/999", nil)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestUpdateAccount(t *testing.T) {
	r, _, db := setupTestRouter(t)
	createTestUser(t, db)

	createResp := executeRequest(r, "POST", "/api/v1/accounts", map[string]any{
		"name": "Checking",
		"type": "checking",
	})
	created := parseResponse(t, createResp)
	id := int(created["data"].(map[string]any)["id"].(float64))

	w := executeRequest(r, "PUT", fmt.Sprintf("/api/v1/accounts/%d", id), map[string]any{
		"name": "Premium Checking",
	})
	assert.Equal(t, http.StatusOK, w.Code)

	updated := parseResponse(t, w)
	assert.Equal(t, "Premium Checking", updated["data"].(map[string]any)["name"])
}

func TestDeleteAccount(t *testing.T) {
	r, _, db := setupTestRouter(t)
	createTestUser(t, db)

	createResp := executeRequest(r, "POST", "/api/v1/accounts", map[string]any{
		"name": "Checking",
		"type": "checking",
	})
	created := parseResponse(t, createResp)
	id := int(created["data"].(map[string]any)["id"].(float64))

	w := executeRequest(r, "DELETE", fmt.Sprintf("/api/v1/accounts/%d", id), nil)
	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestDeleteAccount_NotFound(t *testing.T) {
	r, _, db := setupTestRouter(t)
	createTestUser(t, db)

	w := executeRequest(r, "DELETE", "/api/v1/accounts/999", nil)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

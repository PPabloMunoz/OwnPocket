package handler

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func setupTransactionTestData(t *testing.T, r *gin.Engine, db *gorm.DB) (uint, uint) {
	t.Helper()
	createTestUser(t, db)

	accResp := executeRequest(r, "POST", "/api/v1/accounts", map[string]any{
		"name": "Checking",
		"type": "checking",
	})
	acc := parseResponse(t, accResp)
	fromID := uint(acc["data"].(map[string]any)["id"].(float64))

	accResp2 := executeRequest(r, "POST", "/api/v1/accounts", map[string]any{
		"name": "Savings",
		"type": "savings",
	})
	acc2 := parseResponse(t, accResp2)
	toID := uint(acc2["data"].(map[string]any)["id"].(float64))

	executeRequest(r, "POST", "/api/v1/categories", map[string]any{
		"name": "Salary",
		"type": "income",
	})

	return fromID, toID
}

func TestCreateTransaction_Income(t *testing.T) {
	r, _, db := setupTestRouter(t)
	fromID, _ := setupTransactionTestData(t, r, db)

	w := executeRequest(r, "POST", "/api/v1/transactions", map[string]any{
		"account_id":  fromID,
		"amount":      1000.00,
		"type":        "income",
		"date":        "2026-06-15",
		"description": "Salary",
	})
	assert.Equal(t, http.StatusCreated, w.Code)

	resp := parseResponse(t, w)
	data := resp["data"].(map[string]any)
	assert.Equal(t, "income", data["type"])
	assert.Equal(t, 100000.0, data["amount"].(float64))
}

func TestCreateTransaction_Expense(t *testing.T) {
	r, _, db := setupTestRouter(t)
	fromID, _ := setupTransactionTestData(t, r, db)

	executeRequest(r, "POST", "/api/v1/transactions", map[string]any{
		"account_id":  fromID,
		"amount":      1000.00,
		"type":        "income",
		"date":        "2026-06-15",
		"description": "Salary",
	})

	w := executeRequest(r, "POST", "/api/v1/transactions", map[string]any{
		"account_id":  fromID,
		"amount":      50.00,
		"type":        "expense",
		"date":        "2026-06-15",
		"description": "Groceries",
	})
	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestCreateTransaction_Transfer(t *testing.T) {
	r, _, db := setupTestRouter(t)
	fromID, toID := setupTransactionTestData(t, r, db)

	executeRequest(r, "POST", "/api/v1/transactions", map[string]any{
		"account_id":  fromID,
		"amount":      1000.00,
		"type":        "income",
		"date":        "2026-06-15",
		"description": "Salary",
	})

	w := executeRequest(r, "POST", "/api/v1/transactions", map[string]any{
		"account_id":    fromID,
		"amount":        300.00,
		"type":          "transfer",
		"date":          "2026-06-15",
		"description":   "To savings",
		"to_account_id": toID,
	})
	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestCreateTransaction_TransferMissingToAccount(t *testing.T) {
	r, _, db := setupTestRouter(t)
	fromID, _ := setupTransactionTestData(t, r, db)

	w := executeRequest(r, "POST", "/api/v1/transactions", map[string]any{
		"account_id":  fromID,
		"amount":      100.00,
		"type":        "transfer",
		"date":        "2026-06-15",
		"description": "Bad transfer",
	})
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "to_account_id is required")
}

func TestCreateTransaction_InvalidDate(t *testing.T) {
	r, _, db := setupTestRouter(t)
	fromID, _ := setupTransactionTestData(t, r, db)

	w := executeRequest(r, "POST", "/api/v1/transactions", map[string]any{
		"account_id":  fromID,
		"amount":      100.00,
		"type":        "expense",
		"date":        "invalid-date",
		"description": "Bad date",
	})
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid date format")
}

func TestGetTransactions(t *testing.T) {
	r, _, db := setupTestRouter(t)
	fromID, _ := setupTransactionTestData(t, r, db)

	executeRequest(r, "POST", "/api/v1/transactions", map[string]any{
		"account_id":  fromID,
		"amount":      100.00,
		"type":        "income",
		"date":        "2026-06-15",
		"description": "Salary",
	})

	w := executeRequest(r, "GET", "/api/v1/transactions", nil)
	assert.Equal(t, http.StatusOK, w.Code)

	resp := parseResponse(t, w)
	txs := resp["data"].([]any)
	assert.Len(t, txs, 1)
}

func TestGetTransaction_Success(t *testing.T) {
	r, _, db := setupTestRouter(t)
	fromID, _ := setupTransactionTestData(t, r, db)

	createResp := executeRequest(r, "POST", "/api/v1/transactions", map[string]any{
		"account_id":  fromID,
		"amount":      100.00,
		"type":        "income",
		"date":        "2026-06-15",
		"description": "Salary",
	})
	created := parseResponse(t, createResp)
	id := int(created["data"].(map[string]any)["id"].(float64))

	w := executeRequest(r, "GET", fmt.Sprintf("/api/v1/transactions/%d", id), nil)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetTransaction_NotFound(t *testing.T) {
	r, _, db := setupTestRouter(t)
	_, _ = setupTransactionTestData(t, r, db)

	w := executeRequest(r, "GET", "/api/v1/transactions/999", nil)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestUpdateTransaction(t *testing.T) {
	r, _, db := setupTestRouter(t)
	fromID, _ := setupTransactionTestData(t, r, db)

	createResp := executeRequest(r, "POST", "/api/v1/transactions", map[string]any{
		"account_id":  fromID,
		"amount":      100.00,
		"type":        "expense",
		"date":        "2026-06-15",
		"description": "Groceries",
	})
	created := parseResponse(t, createResp)
	id := int(created["data"].(map[string]any)["id"].(float64))

	w := executeRequest(r, "PUT", fmt.Sprintf("/api/v1/transactions/%d", id), map[string]any{
		"description": "Updated groceries",
	})
	assert.Equal(t, http.StatusOK, w.Code)

	updated := parseResponse(t, w)
	assert.Equal(t, "Updated groceries", updated["data"].(map[string]any)["description"])
}

func TestDeleteTransaction(t *testing.T) {
	r, _, db := setupTestRouter(t)
	fromID, _ := setupTransactionTestData(t, r, db)

	createResp := executeRequest(r, "POST", "/api/v1/transactions", map[string]any{
		"account_id":  fromID,
		"amount":      100.00,
		"type":        "income",
		"date":        "2026-06-15",
		"description": "Salary",
	})
	created := parseResponse(t, createResp)
	id := int(created["data"].(map[string]any)["id"].(float64))

	w := executeRequest(r, "DELETE", fmt.Sprintf("/api/v1/transactions/%d", id), nil)
	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestDeleteTransaction_NotFound(t *testing.T) {
	r, _, db := setupTestRouter(t)
	_, _ = setupTransactionTestData(t, r, db)

	w := executeRequest(r, "DELETE", "/api/v1/transactions/999", nil)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

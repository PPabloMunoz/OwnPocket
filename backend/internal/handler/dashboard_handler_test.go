package handler

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetDashboardSummary(t *testing.T) {
	r, _, db := setupTestRouter(t)
	createTestUser(t, db)

	w := executeRequest(r, "GET", "/api/v1/dashboard/summary", nil)
	assert.Equal(t, http.StatusOK, w.Code)

	resp := parseResponse(t, w)
	data := resp["data"].(map[string]any)
	assert.Equal(t, 0.0, data["total_balance"].(float64))
	assert.Equal(t, 0.0, data["monthly_income"].(float64))
	assert.Equal(t, 0.0, data["monthly_expenses"].(float64))
	assert.Empty(t, data["recent_transactions"])
	assert.Empty(t, data["category_summary"])
	assert.Empty(t, data["budgets"])
}

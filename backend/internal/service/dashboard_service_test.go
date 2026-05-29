package service

import (
	"testing"
	"time"

	"github.com/ppablomunoz/ownpocket/backend/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetDashboardSummary_Empty(t *testing.T) {
	svc, db := setupService(t)
	userID := createTestUser(t, db)

	summary, err := svc.GetDashboardSummary(userID)
	require.NoError(t, err)
	assert.Equal(t, model.Amount(0), summary.TotalBalance)
	assert.Equal(t, model.Amount(0), summary.MonthlyIncome)
	assert.Equal(t, model.Amount(0), summary.MonthlyExpenses)
	assert.Empty(t, summary.RecentTxs)
	assert.Empty(t, summary.Budgets)
	assert.Empty(t, summary.CategorySummary)
}

func TestGetDashboardSummary_WithData(t *testing.T) {
	svc, db := setupService(t)
	userID := createTestUser(t, db)

	account, err := svc.CreateAccount(userID, "Checking", "checking", 1, nil)
	require.NoError(t, err)

	invCategory, err := svc.CreateCategory(userID, "Salary", "income", nil, nil, nil)
	require.NoError(t, err)
	expCategory, err := svc.CreateCategory(userID, "Food", "expense", nil, nil, nil)
	require.NoError(t, err)

	_, err = svc.CreateBudget(userID, expCategory.ID, time.Now().Format("2006-01"), model.NewAmountFromFloat(500.00))
	require.NoError(t, err)

	now := time.Now()
	_, err = svc.CreateTransaction(userID, account.ID, model.NewAmountFromFloat(3000.00), "income", now, "Salary", &invCategory.ID, nil, nil)
	require.NoError(t, err)
	_, err = svc.CreateTransaction(userID, account.ID, model.NewAmountFromFloat(150.00), "expense", now, "Groceries", &expCategory.ID, nil, nil)
	require.NoError(t, err)

	summary, err := svc.GetDashboardSummary(userID)
	require.NoError(t, err)

	assert.Equal(t, model.NewAmountFromFloat(2850.00), summary.TotalBalance)
	assert.Equal(t, model.NewAmountFromFloat(3000.00), summary.MonthlyIncome)
	assert.Equal(t, model.NewAmountFromFloat(150.00), summary.MonthlyExpenses)
	assert.Len(t, summary.RecentTxs, 2)
	assert.Len(t, summary.Budgets, 1)
	assert.Equal(t, model.NewAmountFromFloat(150.00), summary.Budgets[0].Spent)
	assert.Len(t, summary.CategorySummary, 1)
}

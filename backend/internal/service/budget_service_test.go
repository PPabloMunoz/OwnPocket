package service

import (
	"testing"

	"github.com/ppablomunoz/ownpocket/backend/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateBudget_Success(t *testing.T) {
	svc, db := setupService(t)
	userID := createTestUser(t, db)

	category, err := svc.CreateCategory(userID, "Food", "expense", nil, nil, nil)
	require.NoError(t, err)

	budget, err := svc.CreateBudget(userID, category.ID, "2026-06", model.NewAmountFromFloat(500.00))
	require.NoError(t, err)
	assert.NotZero(t, budget.ID)
	assert.Equal(t, category.ID, budget.CategoryID)
	assert.Equal(t, "2026-06", budget.Period)
	assert.Equal(t, model.NewAmountFromFloat(500.00), budget.Amount)
}

func TestGetBudgets(t *testing.T) {
	svc, db := setupService(t)
	userID := createTestUser(t, db)

	category, err := svc.CreateCategory(userID, "Food", "expense", nil, nil, nil)
	require.NoError(t, err)

	_, err = svc.CreateBudget(userID, category.ID, "2026-06", model.NewAmountFromFloat(500.00))
	require.NoError(t, err)
	_, err = svc.CreateBudget(userID, category.ID, "2026-07", model.NewAmountFromFloat(300.00))
	require.NoError(t, err)

	budgets, err := svc.GetBudgets(userID)
	require.NoError(t, err)
	assert.Len(t, budgets, 2)
}

func TestGetBudgets_ScopedByUser(t *testing.T) {
	svc, db := setupService(t)
	userID1 := createTestUser(t, db)
	userID2 := createTestUserWithName(t, db, "other")

	category, err := svc.CreateCategory(userID1, "Food", "expense", nil, nil, nil)
	require.NoError(t, err)

	_, err = svc.CreateBudget(userID1, category.ID, "2026-06", model.NewAmountFromFloat(500.00))
	require.NoError(t, err)

	budgets, err := svc.GetBudgets(userID2)
	require.NoError(t, err)
	assert.Empty(t, budgets)
}

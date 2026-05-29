package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateCategory_Success(t *testing.T) {
	svc, db := setupService(t)
	userID := createTestUser(t, db)

	color := "#ff0000"
	icon := "food"
	category, err := svc.CreateCategory(userID, "Food", "expense", nil, &color, &icon)
	require.NoError(t, err)
	assert.NotZero(t, category.ID)
	assert.Equal(t, "Food", category.Name)
	assert.Equal(t, "expense", category.Type)
}

func TestCreateCategory_WithParent(t *testing.T) {
	svc, db := setupService(t)
	userID := createTestUser(t, db)

	parent, err := svc.CreateCategory(userID, "Food & Drink", "expense", nil, nil, nil)
	require.NoError(t, err)

	child, err := svc.CreateCategory(userID, "Restaurants", "expense", &parent.ID, nil, nil)
	require.NoError(t, err)
	assert.Equal(t, parent.ID, *child.ParentID)
}

func TestCreateCategory_Income(t *testing.T) {
	svc, db := setupService(t)
	userID := createTestUser(t, db)

	category, err := svc.CreateCategory(userID, "Salary", "income", nil, nil, nil)
	require.NoError(t, err)
	assert.Equal(t, "income", category.Type)
}

func TestGetCategories(t *testing.T) {
	svc, db := setupService(t)
	userID := createTestUser(t, db)

	_, err := svc.CreateCategory(userID, "Food", "expense", nil, nil, nil)
	require.NoError(t, err)
	_, err = svc.CreateCategory(userID, "Salary", "income", nil, nil, nil)
	require.NoError(t, err)

	categories, err := svc.GetCategories(userID)
	require.NoError(t, err)
	assert.Len(t, categories, 2)
}

func TestGetCategories_ScopedByUser(t *testing.T) {
	svc, db := setupService(t)
	userID1 := createTestUser(t, db)
	userID2 := createTestUserWithName(t, db, "other")

	_, err := svc.CreateCategory(userID1, "Food", "expense", nil, nil, nil)
	require.NoError(t, err)

	categories, err := svc.GetCategories(userID2)
	require.NoError(t, err)
	assert.Empty(t, categories)
}

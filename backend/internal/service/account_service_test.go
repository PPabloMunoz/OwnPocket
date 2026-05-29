package service

import (
	"testing"

	"github.com/ppablomunoz/ownpocket/backend/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateAccount_Success(t *testing.T) {
	svc, db := setupService(t)
	userID := createTestUser(t, db)

	desc := "My main account"
	account, err := svc.CreateAccount(userID, "Checking", "checking", 1, &desc)
	require.NoError(t, err)
	assert.NotZero(t, account.ID)
	assert.Equal(t, userID, account.UserID)
	assert.Equal(t, "Checking", account.Name)
	assert.Equal(t, "checking", account.Type)
	assert.Equal(t, model.Amount(0), account.Balance)
	assert.True(t, account.IsActive)
}

func TestCreateAccount_DefaultCurrency(t *testing.T) {
	svc, db := setupService(t)
	userID := createTestUser(t, db)

	account, err := svc.CreateAccount(userID, "Savings", "savings", 1, nil)
	require.NoError(t, err)
	assert.NotZero(t, account.ID)
	assert.Equal(t, uint(1), account.CurrencyID)
}

func TestGetAccounts(t *testing.T) {
	svc, db := setupService(t)
	userID := createTestUser(t, db)

	_, err := svc.CreateAccount(userID, "Checking", "checking", 1, nil)
	require.NoError(t, err)
	_, err = svc.CreateAccount(userID, "Savings", "savings", 1, nil)
	require.NoError(t, err)

	accounts, err := svc.GetAccounts(userID)
	require.NoError(t, err)
	assert.Len(t, accounts, 2)
}

func TestGetAccounts_OtherUser(t *testing.T) {
	svc, db := setupService(t)
	userID1 := createTestUser(t, db)
	userID2 := createTestUserWithName(t, db, "other")

	_, err := svc.CreateAccount(userID1, "Checking", "checking", 1, nil)
	require.NoError(t, err)

	accounts, err := svc.GetAccounts(userID2)
	require.NoError(t, err)
	assert.Empty(t, accounts)
}

func TestGetAccount_Success(t *testing.T) {
	svc, db := setupService(t)
	userID := createTestUser(t, db)

	created, err := svc.CreateAccount(userID, "Checking", "checking", 1, nil)
	require.NoError(t, err)

	account, err := svc.GetAccount(userID, created.ID)
	require.NoError(t, err)
	assert.Equal(t, created.ID, account.ID)
	assert.Equal(t, "Checking", account.Name)
}

func TestGetAccount_NotFound(t *testing.T) {
	svc, db := setupService(t)
	userID := createTestUser(t, db)

	_, err := svc.GetAccount(userID, 999)
	assert.Error(t, err)
}

func TestGetAccount_OtherUser(t *testing.T) {
	svc, db := setupService(t)
	userID1 := createTestUser(t, db)
	userID2 := createTestUserWithName(t, db, "other")

	created, err := svc.CreateAccount(userID1, "Checking", "checking", 1, nil)
	require.NoError(t, err)

	_, err = svc.GetAccount(userID2, created.ID)
	assert.Error(t, err)
}

func TestUpdateAccount(t *testing.T) {
	svc, db := setupService(t)
	userID := createTestUser(t, db)

	created, err := svc.CreateAccount(userID, "Checking", "checking", 1, nil)
	require.NoError(t, err)

	updated, err := svc.UpdateAccount(userID, created.ID, map[string]any{"name": "Premium Checking"})
	require.NoError(t, err)
	assert.Equal(t, "Premium Checking", updated.Name)
}

func TestUpdateAccount_NotFound(t *testing.T) {
	svc, db := setupService(t)
	userID := createTestUser(t, db)

	_, err := svc.UpdateAccount(userID, 999, map[string]any{"name": "New Name"})
	assert.Error(t, err)
}

func TestDeleteAccount(t *testing.T) {
	svc, db := setupService(t)
	userID := createTestUser(t, db)

	created, err := svc.CreateAccount(userID, "Checking", "checking", 1, nil)
	require.NoError(t, err)

	err = svc.DeleteAccount(userID, created.ID)
	assert.NoError(t, err)

	_, err = svc.GetAccount(userID, created.ID)
	assert.Error(t, err)
}

func TestDeleteAccount_OtherUser(t *testing.T) {
	svc, db := setupService(t)
	userID1 := createTestUser(t, db)
	userID2 := createTestUserWithName(t, db, "other")

	created, err := svc.CreateAccount(userID1, "Checking", "checking", 1, nil)
	require.NoError(t, err)

	err = svc.DeleteAccount(userID2, created.ID)
	assert.Error(t, err)

	_, err = svc.GetAccount(userID1, created.ID)
	assert.NoError(t, err)
}

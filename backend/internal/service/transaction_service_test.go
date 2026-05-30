package service

import (
	"testing"
	"time"

	"github.com/ppablomunoz/ownpocket/backend/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTransactionTest(t *testing.T) (*Service, uint, uint, uint) {
	t.Helper()
	svc, db := setupService(t)
	userID := createTestUser(t, db)

	account1, err := svc.CreateAccount(userID, "Checking", "checking", 1, nil, nil)
	require.NoError(t, err)

	account2, err := svc.CreateAccount(userID, "Savings", "savings", 1, nil, nil)
	require.NoError(t, err)

	return svc, userID, account1.ID, account2.ID
}

func TestCreateTransaction_Income(t *testing.T) {
	svc, userID, accountID, _ := setupTransactionTest(t)
	date := time.Date(2026, 6, 15, 0, 0, 0, 0, time.UTC)

	tx, err := svc.CreateTransaction(userID, accountID, model.NewAmountFromFloat(1000.00), "income", date, "Salary", nil, nil, nil)
	require.NoError(t, err)
	assert.NotZero(t, tx.ID)
	assert.Equal(t, model.NewAmountFromFloat(1000.00), tx.Amount)
	assert.Equal(t, "income", tx.Type)

	account, err := svc.GetAccount(userID, accountID)
	require.NoError(t, err)
	assert.Equal(t, model.NewAmountFromFloat(1000.00), account.Balance)
}

func TestCreateTransaction_Expense(t *testing.T) {
	svc, userID, accountID, _ := setupTransactionTest(t)
	date := time.Date(2026, 6, 15, 0, 0, 0, 0, time.UTC)

	_, err := svc.CreateTransaction(userID, accountID, model.NewAmountFromFloat(1000.00), "income", date, "Salary", nil, nil, nil)
	require.NoError(t, err)

	tx, err := svc.CreateTransaction(userID, accountID, model.NewAmountFromFloat(50.00), "expense", date, "Groceries", nil, nil, nil)
	require.NoError(t, err)
	assert.Equal(t, "expense", tx.Type)

	account, err := svc.GetAccount(userID, accountID)
	require.NoError(t, err)
	assert.Equal(t, model.NewAmountFromFloat(950.00), account.Balance)
}

func TestCreateTransaction_Transfer(t *testing.T) {
	svc, userID, fromID, toID := setupTransactionTest(t)
	date := time.Date(2026, 6, 15, 0, 0, 0, 0, time.UTC)

	_, err := svc.CreateTransaction(userID, fromID, model.NewAmountFromFloat(1000.00), "income", date, "Salary", nil, nil, nil)
	require.NoError(t, err)

	tx, err := svc.CreateTransaction(userID, fromID, model.NewAmountFromFloat(300.00), "transfer", date, "Transfer to savings", nil, &toID, nil)
	require.NoError(t, err)
	assert.Equal(t, "transfer", tx.Type)
	assert.Nil(t, tx.CategoryID)

	fromAccount, err := svc.GetAccount(userID, fromID)
	require.NoError(t, err)
	assert.Equal(t, model.NewAmountFromFloat(700.00), fromAccount.Balance)

	toAccount, err := svc.GetAccount(userID, toID)
	require.NoError(t, err)
	assert.Equal(t, model.NewAmountFromFloat(300.00), toAccount.Balance)
}

func TestCreateTransaction_TransferWithoutToAccount(t *testing.T) {
	svc, userID, accountID, _ := setupTransactionTest(t)
	date := time.Date(2026, 6, 15, 0, 0, 0, 0, time.UTC)

	_, err := svc.CreateTransaction(userID, accountID, model.NewAmountFromFloat(100.00), "transfer", date, "Missing destination", nil, nil, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "requires a destination account")
}

func TestCreateTransaction_InvalidAccount(t *testing.T) {
	svc, userID, _, _ := setupTransactionTest(t)
	date := time.Date(2026, 6, 15, 0, 0, 0, 0, time.UTC)

	_, err := svc.CreateTransaction(userID, 999, model.NewAmountFromFloat(100.00), "expense", date, "Bad account", nil, nil, nil)
	assert.Error(t, err)
}

func TestGetTransactions(t *testing.T) {
	svc, userID, accountID, _ := setupTransactionTest(t)
	date := time.Date(2026, 6, 15, 0, 0, 0, 0, time.UTC)

	_, err := svc.CreateTransaction(userID, accountID, model.NewAmountFromFloat(100.00), "income", date, "Salary", nil, nil, nil)
	require.NoError(t, err)
	_, err = svc.CreateTransaction(userID, accountID, model.NewAmountFromFloat(25.00), "expense", date, "Food", nil, nil, nil)
	require.NoError(t, err)

	txs, err := svc.GetTransactions(userID)
	require.NoError(t, err)
	assert.Len(t, txs, 2)
}

func TestGetTransactions_ScopedByUser(t *testing.T) {
	svc, db := setupService(t)
	userID1 := createTestUser(t, db)
	userID2 := createTestUserWithName(t, db, "other")
	date := time.Date(2026, 6, 15, 0, 0, 0, 0, time.UTC)

	account1, err := svc.CreateAccount(userID1, "Checking", "checking", 1, nil, nil)
	require.NoError(t, err)

	_, err = svc.CreateTransaction(userID1, account1.ID, model.NewAmountFromFloat(100.00), "income", date, "Salary", nil, nil, nil)
	require.NoError(t, err)

	txs, err := svc.GetTransactions(userID2)
	require.NoError(t, err)
	assert.Empty(t, txs)
}

func TestGetTransaction(t *testing.T) {
	svc, userID, accountID, _ := setupTransactionTest(t)
	date := time.Date(2026, 6, 15, 0, 0, 0, 0, time.UTC)

	created, err := svc.CreateTransaction(userID, accountID, model.NewAmountFromFloat(100.00), "income", date, "Salary", nil, nil, nil)
	require.NoError(t, err)

	fetched, err := svc.GetTransaction(userID, created.ID)
	require.NoError(t, err)
	assert.Equal(t, created.ID, fetched.ID)
	assert.Equal(t, "Salary", fetched.Description)
}

func TestGetTransaction_NotFound(t *testing.T) {
	svc, userID, _, _ := setupTransactionTest(t)

	_, err := svc.GetTransaction(userID, 999)
	assert.Error(t, err)
}

func TestUpdateTransaction_FromExpenseToIncome(t *testing.T) {
	svc, userID, accountID, _ := setupTransactionTest(t)
	date := time.Date(2026, 6, 15, 0, 0, 0, 0, time.UTC)

	_, err := svc.CreateTransaction(userID, accountID, model.NewAmountFromFloat(1000.00), "income", date, "Salary", nil, nil, nil)
	require.NoError(t, err)

	tx, err := svc.CreateTransaction(userID, accountID, model.NewAmountFromFloat(100.00), "expense", date, "Groceries", nil, nil, nil)
	require.NoError(t, err)

	account, _ := svc.GetAccount(userID, accountID)
	assert.Equal(t, model.NewAmountFromFloat(900.00), account.Balance)

	updated, err := svc.UpdateTransaction(userID, tx.ID, map[string]any{
		"type":   "income",
		"amount": float64(model.NewAmountFromFloat(200.00)),
	})
	require.NoError(t, err)
	assert.Equal(t, "income", updated.Type)

	account, _ = svc.GetAccount(userID, accountID)
	assert.Equal(t, model.NewAmountFromFloat(1200.00), account.Balance)
}

func TestUpdateTransaction_NotFound(t *testing.T) {
	svc, userID, _, _ := setupTransactionTest(t)

	_, err := svc.UpdateTransaction(userID, 999, map[string]any{"description": "Updated"})
	assert.Error(t, err)
}

func TestDeleteTransaction(t *testing.T) {
	svc, userID, accountID, _ := setupTransactionTest(t)
	date := time.Date(2026, 6, 15, 0, 0, 0, 0, time.UTC)

	_, err := svc.CreateTransaction(userID, accountID, model.NewAmountFromFloat(1000.00), "income", date, "Salary", nil, nil, nil)
	require.NoError(t, err)

	tx, err := svc.CreateTransaction(userID, accountID, model.NewAmountFromFloat(100.00), "expense", date, "Groceries", nil, nil, nil)
	require.NoError(t, err)

	err = svc.DeleteTransaction(userID, tx.ID)
	require.NoError(t, err)

	account, _ := svc.GetAccount(userID, accountID)
	assert.Equal(t, model.NewAmountFromFloat(1000.00), account.Balance)

	_, err = svc.GetTransaction(userID, tx.ID)
	assert.Error(t, err)
}

func TestDeleteTransaction_NotFound(t *testing.T) {
	svc, userID, _, _ := setupTransactionTest(t)

	err := svc.DeleteTransaction(userID, 999)
	assert.Error(t, err)
}

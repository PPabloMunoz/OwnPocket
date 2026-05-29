package service

import (
	"testing"

	"github.com/ppablomunoz/ownpocket/backend/internal/model"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		SkipDefaultTransaction: false,
	})
	require.NoError(t, err)

	db.Exec("PRAGMA foreign_keys = ON;")

	err = db.AutoMigrate(
		&model.User{},
		&model.Currency{},
		&model.Account{},
		&model.Category{},
		&model.Transaction{},
		&model.Budget{},
		&model.Tag{},
		&model.TransactionTag{},
	)
	require.NoError(t, err)

	seedTestCurrencies(t, db)
	return db
}

func seedTestCurrencies(t *testing.T, db *gorm.DB) {
	t.Helper()

	var count int64
	db.Model(&model.Currency{}).Count(&count)
	if count > 0 {
		return
	}

	currencies := []model.Currency{
		{Code: "EUR", Name: "Euro", Symbol: "€", DecimalPlaces: 2},
		{Code: "USD", Name: "US Dollar", Symbol: "$", DecimalPlaces: 2},
	}
	err := db.Create(&currencies).Error
	require.NoError(t, err)
}

func setupService(t *testing.T) (*Service, *gorm.DB) {
	t.Helper()
	db := setupTestDB(t)
	return NewService(db), db
}

func createTestUser(t *testing.T, db *gorm.DB) uint {
	t.Helper()
	return createTestUserWithName(t, db, "testuser")
}

func createTestUserWithName(t *testing.T, db *gorm.DB, username string) uint {
	t.Helper()
	user := model.User{Username: username, PasswordHash: "hash"}
	err := db.Create(&user).Error
	require.NoError(t, err)
	return user.ID
}

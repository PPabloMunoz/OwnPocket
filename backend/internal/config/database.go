package config

import (
	"github.com/ppablomunoz/ownpocket/backend/internal/model"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func NewDatabase(dbPath string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Enable foreign keys
	db.Exec("PRAGMA foreign_keys = ON;")

	// Auto migrate (use in development)
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
	if err != nil {
		return nil, err
	}

	return db, nil
}

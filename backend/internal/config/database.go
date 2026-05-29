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

	if err := seedCurrencies(db); err != nil {
		return nil, err
	}

	return db, nil
}

func seedCurrencies(db *gorm.DB) error {
	var count int64
	db.Model(&model.Currency{}).Count(&count)
	if count > 0 {
		return nil
	}

	currencies := []model.Currency{
		{Code: "EUR", Name: "Euro", Symbol: "€", DecimalPlaces: 2},
		{Code: "USD", Name: "US Dollar", Symbol: "$", DecimalPlaces: 2},
		{Code: "GBP", Name: "British Pound", Symbol: "£", DecimalPlaces: 2},
		{Code: "JPY", Name: "Japanese Yen", Symbol: "¥", DecimalPlaces: 0},
		{Code: "CHF", Name: "Swiss Franc", Symbol: "CHF", DecimalPlaces: 2},
		{Code: "CAD", Name: "Canadian Dollar", Symbol: "CA$", DecimalPlaces: 2},
		{Code: "AUD", Name: "Australian Dollar", Symbol: "A$", DecimalPlaces: 2},
		{Code: "CNY", Name: "Chinese Yuan", Symbol: "¥", DecimalPlaces: 2},
		{Code: "SEK", Name: "Swedish Krona", Symbol: "kr", DecimalPlaces: 2},
		{Code: "NOK", Name: "Norwegian Krone", Symbol: "kr", DecimalPlaces: 2},
	}

	return db.Create(&currencies).Error
}

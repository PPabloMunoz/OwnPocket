package service

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/ppablomunoz/ownpocket/backend/internal/model"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func (s *Service) RegisterUser(username, password, email string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := model.User{
		Username:     username,
		PasswordHash: string(hashedPassword),
		Email:        &email,
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&user).Error; err != nil {
			return err
		}

		// Seed default categories
		return seedDefaultCategories(tx, user.ID)
	})
}

func seedDefaultCategories(db *gorm.DB, userID uint) error {
	categories := []model.Category{
		{UserID: userID, Name: "Housing", Type: "expense", Color: ptr("#ef4444")},
		{UserID: userID, Name: "Groceries", Type: "expense", Color: ptr("#f97316")},
		{UserID: userID, Name: "Transport", Type: "expense", Color: ptr("#eab308")},
		{UserID: userID, Name: "Dining Out", Type: "expense", Color: ptr("#84cc16")},
		{UserID: userID, Name: "Utilities", Type: "expense", Color: ptr("#10b981")},
		{UserID: userID, Name: "Health", Type: "expense", Color: ptr("#06b6d4")},
		{UserID: userID, Name: "Shopping", Type: "expense", Color: ptr("#3b82f6")},
		{UserID: userID, Name: "Entertainment", Type: "expense", Color: ptr("#6366f1")},
		{UserID: userID, Name: "Salary", Type: "income", Color: ptr("#8b5cf6")},
		{UserID: userID, Name: "Gifts", Type: "income", Color: ptr("#d946ef")},
		{UserID: userID, Name: "Investment", Type: "income", Color: ptr("#f43f5e")},
	}
	return db.Create(&categories).Error
}

func ptr[T any](v T) *T {
	return &v
}

func (s *Service) Login(username, password, jwtSecret string) (string, error) {
	var user model.User
	if err := s.db.Where("username = ?", username).First(&user).Error; err != nil {
		return "", errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", errors.New("invalid credentials")
	}

	// Generate JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(24 * 30 * time.Hour).Unix(), // 30 days
	})

	return token.SignedString([]byte(jwtSecret))
}

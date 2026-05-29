package service

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/ppablomunoz/ownpocket/backend/internal/model"
	"golang.org/x/crypto/bcrypt"
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

	return s.db.Create(&user).Error
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

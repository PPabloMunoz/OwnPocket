package service

import (
	"gorm.io/gorm"
)

// Service holds all business logic
type Service struct {
	db *gorm.DB
}

// NewService creates a new service instance
func NewService(db *gorm.DB) *Service {
	return &Service{db: db}
}

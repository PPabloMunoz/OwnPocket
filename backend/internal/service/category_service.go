package service

import (
	"github.com/ppablomunoz/ownpocket/backend/internal/model"
)

func (s *Service) CreateCategory(userID uint, name, catType string, parentID *uint, color, icon *string) (*model.Category, error) {
	category := model.Category{
		UserID:   userID,
		Name:     name,
		Type:     catType,
		ParentID: parentID,
		Color:    color,
		Icon:     icon,
	}
	if err := s.db.Create(&category).Error; err != nil {
		return nil, err
	}
	return &category, nil
}

func (s *Service) GetCategories(userID uint) ([]model.Category, error) {
	var categories []model.Category
	if err := s.db.Where("user_id = ?", userID).Preload("Parent").Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

package service

import (
	"errors"

	"github.com/ppablomunoz/ownpocket/backend/internal/model"
	"gorm.io/gorm"
)

var ErrCategoryHasDependencies = errors.New("category has child categories or transactions")

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

func (s *Service) UpdateCategory(userID, categoryID uint, updates map[string]any) (*model.Category, error) {
	var category model.Category
	if err := s.db.Where("user_id = ? AND id = ?", userID, categoryID).First(&category).Error; err != nil {
		return nil, err
	}
	if err := s.db.Model(&category).Updates(updates).Error; err != nil {
		return nil, err
	}
	s.db.Where("user_id = ? AND id = ?", userID, categoryID).Preload("Parent").First(&category)
	return &category, nil
}

func (s *Service) DeleteCategory(userID, categoryID uint) error {
	var count int64
	s.db.Model(&model.Category{}).Where("parent_id = ? AND user_id = ?", categoryID, userID).Count(&count)
	if count > 0 {
		return ErrCategoryHasDependencies
	}
	s.db.Model(&model.Transaction{}).Where("category_id = ? AND user_id = ?", categoryID, userID).Count(&count)
	if count > 0 {
		return ErrCategoryHasDependencies
	}
	result := s.db.Where("user_id = ? AND id = ?", userID, categoryID).Delete(&model.Category{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

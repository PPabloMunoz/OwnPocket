package service

import (
	"github.com/ppablomunoz/ownpocket/backend/internal/model"
)

func (s *Service) CreateBudget(userID, categoryID uint, period string, amount model.Amount) (*model.Budget, error) {
	budget := model.Budget{
		UserID:     userID,
		CategoryID: categoryID,
		Period:     period,
		Amount:     amount,
	}
	if err := s.db.Create(&budget).Error; err != nil {
		return nil, err
	}
	return &budget, nil
}

func (s *Service) GetBudgets(userID uint) ([]model.Budget, error) {
	var budgets []model.Budget
	if err := s.db.Where("user_id = ?", userID).Preload("Category").Find(&budgets).Error; err != nil {
		return nil, err
	}
	return budgets, nil
}

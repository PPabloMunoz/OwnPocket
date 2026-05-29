package service

import (
	"github.com/ppablomunoz/ownpocket/backend/internal/model"
	"gorm.io/gorm"
)

func (s *Service) CreateAccount(
	userID uint,
	name, accountType string,
	currencyID uint,
	description *string,
) (*model.Account, error) {
	account := model.Account{
		UserID:      userID,
		Name:        name,
		Type:        accountType,
		Balance:     0,
		CurrencyID:  currencyID,
		Description: description,
		IsActive:    true,
	}
	if err := s.db.Create(&account).Error; err != nil {
		return nil, err
	}
	return &account, nil
}

func (s *Service) GetAccounts(userID uint) ([]model.Account, error) {
	var accounts []model.Account
	if err := s.db.Where("user_id = ?", userID).Preload("Currency").Find(&accounts).Error; err != nil {
		return nil, err
	}
	return accounts, nil
}

func (s *Service) GetAccount(userID, accountID uint) (*model.Account, error) {
	var account model.Account
	if err := s.db.Where("user_id = ? AND id = ?", userID, accountID).Preload("Currency").First(&account).Error; err != nil {
		return nil, err
	}
	return &account, nil
}

func (s *Service) UpdateAccount(userID, accountID uint, updates map[string]any) (*model.Account, error) {
	var account model.Account
	if err := s.db.Where("user_id = ? AND id = ?", userID, accountID).First(&account).Error; err != nil {
		return nil, err
	}
	if err := s.db.Model(&account).Updates(updates).Error; err != nil {
		return nil, err
	}
	s.db.First(&account)
	return &account, nil
}

func (s *Service) DeleteAccount(userID, accountID uint) error {
	result := s.db.Where("user_id = ? AND id = ?", userID, accountID).Delete(&model.Account{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

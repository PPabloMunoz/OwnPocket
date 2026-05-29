package service

import (
	"time"

	"github.com/ppablomunoz/ownpocket/backend/internal/model"
)

func (s *Service) CreateTransaction(userID, accountID uint, amount model.Amount, txType string, date time.Time, description string, categoryID, toAccountID *uint, notes *string) (*model.Transaction, error) {
	tx := model.Transaction{
		UserID:      userID,
		AccountID:   accountID,
		Amount:      amount,
		Type:        txType,
		Date:        date,
		Description: description,
		CategoryID:  categoryID,
		ToAccountID: toAccountID,
		Notes:       notes,
	}
	if err := s.db.Create(&tx).Error; err != nil {
		return nil, err
	}
	return &tx, nil
}

func (s *Service) GetTransactions(userID uint) ([]model.Transaction, error) {
	var txs []model.Transaction
	if err := s.db.Where("user_id = ?", userID).
		Preload("Account").
		Preload("ToAccount").
		Preload("Category").
		Preload("Tags").
		Order("date DESC").
		Find(&txs).Error; err != nil {
		return nil, err
	}
	return txs, nil
}

func (s *Service) GetTransaction(userID, txID uint) (*model.Transaction, error) {
	var tx model.Transaction
	if err := s.db.Where("user_id = ? AND id = ?", userID, txID).
		Preload("Account").
		Preload("ToAccount").
		Preload("Category").
		Preload("Tags").
		First(&tx).Error; err != nil {
		return nil, err
	}
	return &tx, nil
}

func (s *Service) UpdateTransaction(userID, txID uint, updates map[string]any) (*model.Transaction, error) {
	var tx model.Transaction
	if err := s.db.Where("user_id = ? AND id = ?", userID, txID).First(&tx).Error; err != nil {
		return nil, err
	}
	if err := s.db.Model(&tx).Updates(updates).Error; err != nil {
		return nil, err
	}
	s.db.First(&tx)
	return &tx, nil
}

func (s *Service) DeleteTransaction(userID, txID uint) error {
	result := s.db.Where("user_id = ? AND id = ?", userID, txID).Delete(&model.Transaction{})
	return result.Error
}

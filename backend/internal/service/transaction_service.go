package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/ppablomunoz/ownpocket/backend/internal/model"
	"gorm.io/gorm"
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

	if txType == "transfer" {
		tx.CategoryID = nil
	}

	err := s.db.Transaction(func(dbTx *gorm.DB) error {
		if err := dbTx.Create(&tx).Error; err != nil {
			return err
		}
		return updateBalances(dbTx, userID, accountID, toAccountID, amount, txType, 1)
	})
	if err != nil {
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
		Order("date DESC, id DESC").
		Find(&txs).Error; err != nil {
		return nil, err
	}
	return txs, nil
}

func (s *Service) GetTransactionsPaginated(userID uint, page, pageSize int) ([]model.Transaction, int64, error) {
	var total int64
	if err := s.db.Model(&model.Transaction{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var txs []model.Transaction
	offset := (page - 1) * pageSize
	if err := s.db.Where("user_id = ?", userID).
		Preload("Account").
		Preload("ToAccount").
		Preload("Category").
		Preload("Tags").
		Order("date DESC, id DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&txs).Error; err != nil {
		return nil, 0, err
	}
	return txs, total, nil
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

	err := s.db.Transaction(func(dbTx *gorm.DB) error {
		if err := dbTx.Where("user_id = ? AND id = ?", userID, txID).First(&tx).Error; err != nil {
			return err
		}

		// Reverse original balance effect
		if err := updateBalances(dbTx, userID, tx.AccountID, tx.ToAccountID, tx.Amount, tx.Type, -1); err != nil {
			return err
		}

		// Apply updates
		if err := dbTx.Model(&tx).Updates(updates).Error; err != nil {
			return err
		}

		// Re-read updated tx
		if err := dbTx.First(&tx).Error; err != nil {
			return err
		}

		// Apply new balance effect
		return updateBalances(dbTx, userID, tx.AccountID, tx.ToAccountID, tx.Amount, tx.Type, 1)
	})
	if err != nil {
		return nil, err
	}
	return &tx, nil
}

func (s *Service) DeleteTransaction(userID, txID uint) error {
	return s.db.Transaction(func(dbTx *gorm.DB) error {
		var tx model.Transaction
		if err := dbTx.Where("user_id = ? AND id = ?", userID, txID).First(&tx).Error; err != nil {
			return err
		}

		// Reverse balance effect
		if err := updateBalances(dbTx, userID, tx.AccountID, tx.ToAccountID, tx.Amount, tx.Type, -1); err != nil {
			return err
		}

		if err := dbTx.Delete(&tx).Error; err != nil {
			return err
		}
		return nil
	})
}

// updateBalances adjusts account balances atomically within a transaction.
// sign: +1 to apply, -1 to reverse.
func updateBalances(dbTx *gorm.DB, userID, accountID uint, toAccountID *uint, amount model.Amount, txType string, sign int) error {
	switch txType {
	case "income":
		result := dbTx.Model(&model.Account{}).
			Where("id = ? AND user_id = ?", accountID, userID).
			Update("balance", gorm.Expr("balance + ?", sign*int(amount.Cents())))
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return errors.New("account not found")
		}

	case "expense":
		result := dbTx.Model(&model.Account{}).
			Where("id = ? AND user_id = ?", accountID, userID).
			Update("balance", gorm.Expr("balance - ?", sign*int(amount.Cents())))
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return errors.New("account not found")
		}

	case "transfer":
		if toAccountID == nil {
			return fmt.Errorf("transfer requires a destination account")
		}
		result := dbTx.Model(&model.Account{}).
			Where("id = ? AND user_id = ?", accountID, userID).
			Update("balance", gorm.Expr("balance - ?", sign*int(amount.Cents())))
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return errors.New("source account not found")
		}

		result = dbTx.Model(&model.Account{}).
			Where("id = ? AND user_id = ?", *toAccountID, userID).
			Update("balance", gorm.Expr("balance + ?", sign*int(amount.Cents())))
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return errors.New("destination account not found")
		}
	}
	return nil
}

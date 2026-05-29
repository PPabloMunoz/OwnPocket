package service

import "github.com/ppablomunoz/ownpocket/backend/internal/model"

type DashboardSummary struct {
	TotalBalance    model.Amount            `json:"total_balance"`
	MonthlyIncome   model.Amount            `json:"monthly_income"`
	MonthlyExpenses model.Amount            `json:"monthly_expenses"`
	RecentTxs       []model.Transaction     `json:"recent_transactions"`
	CategorySummary []CategorySummaryItem   `json:"category_summary"`
	Budgets         []BudgetWithSpent       `json:"budgets"`
}

type CategorySummaryItem struct {
	Category model.Category `json:"category"`
	Total    model.Amount   `json:"total"`
}

type BudgetWithSpent struct {
	model.Budget
	Spent model.Amount `json:"spent"`
}

func (s *Service) GetDashboardSummary(userID uint) (*DashboardSummary, error) {
	var summary DashboardSummary

	var accounts []model.Account
	if err := s.db.Where("user_id = ?", userID).Find(&accounts).Error; err != nil {
		return nil, err
	}
	for _, a := range accounts {
		summary.TotalBalance += a.Balance
	}

	var recentTxs []model.Transaction
	if err := s.db.Where("user_id = ?", userID).
		Preload("Account").
		Preload("Category").
		Order("date DESC").
		Limit(10).
		Find(&recentTxs).Error; err != nil {
		return nil, err
	}
	summary.RecentTxs = recentTxs

	s.db.Model(&model.Transaction{}).
		Where("user_id = ? AND type = 'income' AND strftime('%Y-%m', date) = strftime('%Y-%m', 'now')", userID).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&summary.MonthlyIncome)

	s.db.Model(&model.Transaction{}).
		Where("user_id = ? AND type = 'expense' AND strftime('%Y-%m', date) = strftime('%Y-%m', 'now')", userID).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&summary.MonthlyExpenses)

	var budgets []model.Budget
	if err := s.db.Where("user_id = ?", userID).Preload("Category").Find(&budgets).Error; err != nil {
		return nil, err
	}
	for _, b := range budgets {
		var spent model.Amount
		s.db.Model(&model.Transaction{}).
			Where("user_id = ? AND category_id = ? AND strftime('%Y-%m', date) = ? AND type = 'expense'", userID, b.CategoryID, b.Period).
			Select("COALESCE(SUM(amount), 0)").
			Scan(&spent)
		summary.Budgets = append(summary.Budgets, BudgetWithSpent{Budget: b, Spent: spent})
	}

	type catRow struct {
		CategoryID uint
		Total      model.Amount
	}
	var rows []catRow
	s.db.Model(&model.Transaction{}).
		Select("category_id, SUM(amount) as total").
		Where("user_id = ? AND type = 'expense'", userID).
		Group("category_id").
		Scan(&rows)
	for _, r := range rows {
		var cat model.Category
		if err := s.db.First(&cat, r.CategoryID).Error; err != nil {
			continue
		}
		summary.CategorySummary = append(summary.CategorySummary, CategorySummaryItem{
			Category: cat,
			Total:    r.Total,
		})
	}

	return &summary, nil
}

import type { Transaction } from "./transaction";
import type { Category } from "./category";
import type { Budget } from "./budget";

export interface CategorySummaryItem {
  category: Category;
  total: number;
}

export interface BudgetWithSpent extends Budget {
  spent: number;
}

export interface DashboardSummary {
  total_balance: number;
  monthly_income: number;
  monthly_expenses: number;
  recent_transactions: Transaction[];
  category_summary: CategorySummaryItem[];
  budgets: BudgetWithSpent[];
}

import type { Category } from "./category";

export interface Budget {
  id: number;
  user_id: number;
  category_id: number;
  period: string;
  amount: number;
  created_at: string;
  category?: Category;
}

export interface CreateBudgetRequest {
  category_id: number;
  period: string;
  amount: number;
}

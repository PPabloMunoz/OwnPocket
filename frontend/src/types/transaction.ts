import type { Account } from "./account";
import type { Category } from "./category";

export interface Transaction {
  id: number;
  user_id: number;
  account_id: number;
  to_account_id: number | null;
  category_id: number | null;
  amount: number;
  type: "income" | "expense" | "transfer";
  date: string;
  description: string;
  notes: string | null;
  reconciled: boolean;
  created_at: string;
  updated_at: string;
  account?: Account;
  to_account?: Account;
  category?: Category;
  tags?: Tag[];
}

export interface CreateTransactionRequest {
  account_id: number;
  to_account_id?: number;
  category_id?: number;
  amount: number;
  type: Transaction["type"];
  date: string;
  description: string;
  notes?: string;
  tag_ids?: number[];
}

export interface Tag {
  id: number;
  user_id: number;
  name: string;
  color: string | null;
  created_at: string;
}

export interface Account {
  id: number;
  user_id: number;
  name: string;
  type: "checking" | "savings" | "credit_card" | "cash" | "investment" | "loan";
  balance: number;
  currency_id: number;
  description: string | null;
  is_active: boolean;
  created_at: string;
  updated_at: string;
  currency?: Currency;
}

export interface CreateAccountRequest {
  name: string;
  type: Account["type"];
  balance?: number;
  currency_id?: number;
  description?: string;
}

export interface UpdateAccountRequest {
  name?: string;
  type?: Account["type"];
  description?: string;
  is_active?: boolean;
}

export interface Currency {
  id: number;
  code: string;
  name: string;
  symbol: string;
  decimal_places: number;
}

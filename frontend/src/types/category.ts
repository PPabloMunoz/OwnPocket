export interface Category {
  id: number;
  user_id: number;
  name: string;
  parent_id: number | null;
  color: string | null;
  icon: string | null;
  type: "income" | "expense";
  created_at: string;
  parent?: Category;
}

export interface CreateCategoryRequest {
  name: string;
  type: Category["type"];
  parent_id?: number;
  color?: string;
  icon?: string;
}

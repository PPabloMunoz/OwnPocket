export const queryKeys = {
  auth: {
    me: ["auth", "me"] as const,
  },
  accounts: {
    all: ["accounts"] as const,
    detail: (id: number) => ["accounts", id] as const,
  },
  transactions: {
    all: ["transactions"] as const,
    detail: (id: number) => ["transactions", id] as const,
  },
  categories: {
    all: ["categories"] as const,
  },
  budgets: {
    all: ["budgets"] as const,
  },
  dashboard: {
    summary: ["dashboard", "summary"] as const,
  },
};

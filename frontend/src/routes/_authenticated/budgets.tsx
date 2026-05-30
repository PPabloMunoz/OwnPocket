import { createFileRoute } from "@tanstack/react-router";
import { useState, type FormEvent } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { Plus, PiggyBank, X, AlertCircle } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Card } from "@/components/ui/card";
import { api } from "@/lib/api";
import { queryKeys } from "@/lib/query-keys";
import { formatCents } from "@/lib/utils";
import type { Budget, CreateBudgetRequest } from "@/types/budget";
import type { Category } from "@/types/category";

export const Route = createFileRoute("/_authenticated/budgets")({
  component: BudgetsPage,
});

function groupByPeriod(budgets: Budget[]): Record<string, Budget[]> {
  const groups: Record<string, Budget[]> = {};
  for (const budget of budgets) {
    (groups[budget.period] ??= []).push(budget);
  }
  return groups;
}

function formatPeriod(period: string): string {
  const [year, month] = period.split("-");
  return new Date(Number(year), Number(month) - 1).toLocaleDateString("en-US", {
    year: "numeric",
    month: "long",
  });
}

function BudgetsPage() {
  const queryClient = useQueryClient();
  const [showForm, setShowForm] = useState(false);
  const [categoryId, setCategoryId] = useState<number | "">("");
  const [period, setPeriod] = useState("");
  const [amount, setAmount] = useState("");

  const { data: budgets = [], isLoading } = useQuery({
    queryKey: queryKeys.budgets.all,
    queryFn: () => api.get<Budget[]>("/budgets"),
  });

  const { data: categories = [] } = useQuery({
    queryKey: queryKeys.categories.all,
    queryFn: () => api.get<Category[]>("/categories"),
  });

  const createMutation = useMutation({
    mutationFn: (body: CreateBudgetRequest) => api.post<Budget>("/budgets", body),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.budgets.all });
      resetForm();
    },
  });

  function resetForm() {
    setShowForm(false);
    setCategoryId("");
    setPeriod("");
    setAmount("");
  }

  const handleSubmit = (e: FormEvent) => {
    e.preventDefault();
    if (!categoryId || !period.trim() || !amount.trim()) return;
    createMutation.mutate({
      category_id: Number(categoryId),
      period: period.trim(),
      amount: parseFloat(amount),
    });
  };

  const isPending = createMutation.isPending;
  const mutationError = createMutation.error;
  const grouped = groupByPeriod(budgets);
  const sortedPeriods = Object.keys(grouped).sort().reverse();
  const expenseCategories = categories.filter((c) => c.type === "expense");

  return (
    <div className="space-y-8">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-zinc-900 dark:text-zinc-50">Budgets</h1>
          <p className="mt-1 text-zinc-500 dark:text-zinc-400">
            Set monthly spending limits by category.
          </p>
        </div>
        {!showForm && (
          <Button onClick={() => setShowForm(true)}>
            <Plus className="mr-1.5 h-4 w-4" />
            Add budget
          </Button>
        )}
      </div>

      {mutationError && (
        <div className="flex items-start gap-3 rounded-xl border border-red-200 bg-red-50/80 px-4 py-3 backdrop-blur-xl dark:border-red-900/50 dark:bg-red-950/50">
          <AlertCircle className="mt-0.5 h-4 w-4 shrink-0 text-red-500 dark:text-red-400" />
          <p className="flex-1 text-sm text-red-700 dark:text-red-300">
            {mutationError instanceof Error ? mutationError.message : "An error occurred"}
          </p>
          <button
            onClick={() => createMutation.reset()}
            className="rounded-lg p-1 text-red-500 transition-colors hover:bg-red-100 dark:text-red-400 dark:hover:bg-red-900/50"
            aria-label="Dismiss"
          >
            <X className="h-3.5 w-3.5" />
          </button>
        </div>
      )}

      {showForm && (
        <Card>
          <h2 className="mb-4 text-lg font-semibold text-zinc-900 dark:text-zinc-50">New budget</h2>
          <form onSubmit={handleSubmit} className="flex flex-col gap-4">
            <div className="grid gap-4 sm:grid-cols-2">
              <div className="space-y-1.5">
                <label
                  htmlFor="category"
                  className="block text-sm font-medium text-zinc-700 dark:text-zinc-300"
                >
                  Category
                </label>
                <select
                  id="category"
                  value={categoryId}
                  onChange={(e) => setCategoryId(e.target.value ? Number(e.target.value) : "")}
                  required
                  className="block w-full rounded-xl border border-zinc-200 bg-white px-3 py-2 text-sm text-zinc-900 focus:border-zinc-400 focus:outline-none dark:border-zinc-700 dark:bg-zinc-800 dark:text-zinc-100 dark:focus:border-zinc-500"
                >
                  <option value="">Select a category</option>
                  {expenseCategories.map((cat) => (
                    <option key={cat.id} value={cat.id}>
                      {cat.name}
                    </option>
                  ))}
                </select>
              </div>
              <Input
                label="Period"
                id="period"
                value={period}
                onChange={(e) => setPeriod(e.target.value)}
                required
                placeholder="YYYY-MM"
                pattern="\d{4}-\d{2}"
              />
              <Input
                label="Monthly limit"
                id="amount"
                type="number"
                step="0.01"
                min="0"
                value={amount}
                onChange={(e) => setAmount(e.target.value)}
                required
                placeholder="0.00"
              />
            </div>
            <div className="flex justify-end gap-3">
              <Button type="button" variant="ghost" onClick={resetForm}>
                Cancel
              </Button>
              <Button
                type="submit"
                disabled={isPending || !categoryId || !period.trim() || !amount.trim()}
              >
                {isPending ? "Saving..." : "Create"}
              </Button>
            </div>
          </form>
        </Card>
      )}

      {isLoading ? (
        <Card>
          <p className="text-center text-sm text-zinc-400 dark:text-zinc-500">Loading budgets...</p>
        </Card>
      ) : budgets.length === 0 ? (
        <Card>
          <div className="flex flex-col items-center gap-3 py-12 text-center">
            <PiggyBank className="h-10 w-10 text-zinc-300 dark:text-zinc-600" />
            <p className="text-sm text-zinc-400 dark:text-zinc-500">
              No budgets yet. Create a budget to track your spending.
            </p>
          </div>
        </Card>
      ) : (
        <div className="space-y-8">
          {sortedPeriods.map((period) => (
            <div key={period}>
              <h2 className="mb-3 text-sm font-semibold uppercase tracking-wider text-zinc-500 dark:text-zinc-400">
                {formatPeriod(period)}
              </h2>
              <div className="space-y-1">
                {grouped[period].map((budget) => (
                  <div
                    key={budget.id}
                    className="flex items-center justify-between rounded-xl border border-zinc-200 bg-white/80 px-4 py-3 shadow-sm backdrop-blur-xl dark:border-zinc-800 dark:bg-zinc-900/80"
                  >
                    <div className="flex items-center gap-3">
                      {budget.category && (
                        <span
                          className="h-3 w-3 shrink-0 rounded-full"
                          style={{
                            backgroundColor: budget.category.color ?? "oklch(0.5 0 0)",
                          }}
                        />
                      )}
                      <span className="text-sm font-medium text-zinc-900 dark:text-zinc-50">
                        {budget.category?.name ?? "Unknown"}
                      </span>
                    </div>
                    <span className="text-sm font-semibold text-zinc-900 dark:text-zinc-50">
                      {formatCents(budget.amount)}
                    </span>
                  </div>
                ))}
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}

import { createFileRoute } from "@tanstack/react-router";
import { useState, type FormEvent } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import {
  Plus,
  ArrowLeftRight,
  ArrowDownRight,
  ArrowUpRight,
  ArrowRightLeft,
  X,
  AlertCircle,
  Pencil,
  Trash2,
} from "lucide-react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Card } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { api } from "@/lib/api";
import { queryKeys } from "@/lib/query-keys";
import { formatCents } from "@/lib/utils";
import type { Transaction } from "@/types/transaction";
import type { Account } from "@/types/account";
import type { Category } from "@/types/category";
import type { PaginatedData } from "@/types/api";

export const Route = createFileRoute("/_authenticated/transactions")({
  component: TransactionsPage,
});

const PAGE_SIZE = 20;

const TXN_TYPES = ["income", "expense", "transfer"] as const;

function todayStr(): string {
  return new Date().toISOString().slice(0, 10);
}

function groupByDate(transactions: Transaction[]): Record<string, Transaction[]> {
  const groups: Record<string, Transaction[]> = {};
  for (const tx of transactions) {
    (groups[tx.date] ??= []).push(tx);
  }
  return groups;
}

function formatDate(dateStr: string): string {
  const d = new Date(dateStr.slice(0, 10) + "T00:00:00");
  return d.toLocaleDateString("en-US", {
    weekday: "short",
    month: "short",
    day: "numeric",
    year: "numeric",
  });
}

function TransactionsPage() {
  const queryClient = useQueryClient();
  const [showForm, setShowForm] = useState(false);
  const [editingId, setEditingId] = useState<number | null>(null);
  const [txnType, setTxnType] = useState<string>("expense");
  const [accountId, setAccountId] = useState<number | "">("");
  const [toAccountId, setToAccountId] = useState<number | "">("");
  const [categoryId, setCategoryId] = useState<number | "">("");
  const [amount, setAmount] = useState("");
  const [date, setDate] = useState(todayStr());
  const [description, setDescription] = useState("");
  const [page, setPage] = useState(1);

  const { data, isLoading } = useQuery({
    queryKey: queryKeys.transactions.paginated(page, PAGE_SIZE),
    queryFn: () =>
      api.get<PaginatedData<Transaction>>(`/transactions?page=${page}&page_size=${PAGE_SIZE}`),
  });

  const transactions = data?.items ?? [];
  const totalPages = data?.total_pages ?? 1;
  const total = data?.total ?? 0;

  const { data: accounts = [] } = useQuery({
    queryKey: queryKeys.accounts.all,
    queryFn: () => api.get<Account[]>("/accounts"),
  });

  const { data: categories = [] } = useQuery({
    queryKey: queryKeys.categories.all,
    queryFn: () => api.get<Category[]>("/categories"),
  });

  function refetchAll() {
    queryClient.invalidateQueries({ queryKey: queryKeys.transactions.all });
    queryClient.invalidateQueries({ queryKey: queryKeys.accounts.all });
    setPage(1);
  }

  const createMutation = useMutation({
    mutationFn: (body: {
      account_id: number;
      to_account_id?: number;
      category_id?: number;
      amount: number;
      type: string;
      date: string;
      description: string;
    }) => api.post<Transaction>("/transactions", body),
    onSuccess: () => {
      refetchAll();
      resetForm();
    },
  });

  const updateMutation = useMutation({
    mutationFn: ({
      id,
      ...body
    }: {
      id: number;
      account_id?: number;
      to_account_id?: number;
      category_id?: number;
      amount?: number;
      type?: string;
      date?: string;
      description?: string;
    }) => api.put<Transaction>(`/transactions/${id}`, body),
    onSuccess: () => {
      refetchAll();
      resetForm();
    },
  });

  const deleteMutation = useMutation({
    mutationFn: (id: number) => api.delete(`/transactions/${id}`),
    onSuccess: () => {
      refetchAll();
    },
  });

  function resetForm() {
    setShowForm(false);
    setEditingId(null);
    setTxnType("expense");
    setAccountId("");
    setToAccountId("");
    setCategoryId("");
    setAmount("");
    setDate(todayStr());
    setDescription("");
  }

  function populateForm(tx: Transaction) {
    setEditingId(tx.id);
    setTxnType(tx.type);
    setAccountId(tx.account_id);
    setToAccountId(tx.to_account_id ?? "");
    setCategoryId(tx.category_id ?? "");
    setAmount((tx.amount / 100).toFixed(2));
    setDate(tx.date.slice(0, 10));
    setDescription(tx.description);
    setShowForm(true);
    window.scrollTo({ top: 0, behavior: "smooth" });
  }

  const handleSubmit = (e: FormEvent) => {
    e.preventDefault();
    if (!accountId || !amount.trim() || !date || !description.trim()) return;

    const payload: {
      account_id: number;
      to_account_id?: number;
      category_id?: number;
      amount: number;
      type: string;
      date: string;
      description: string;
    } = {
      account_id: Number(accountId),
      amount: parseFloat(amount),
      type: txnType,
      date,
      description: description.trim(),
    };

    if (txnType === "transfer" && toAccountId) {
      payload.to_account_id = Number(toAccountId);
    }
    if (txnType !== "transfer" && categoryId) {
      payload.category_id = Number(categoryId);
    }

    if (editingId) {
      updateMutation.mutate({ id: editingId, ...payload });
    } else {
      createMutation.mutate(payload);
    }
  };

  const isPending = createMutation.isPending || updateMutation.isPending;
  const mutationError = createMutation.error ?? updateMutation.error ?? deleteMutation.error;

  const grouped = groupByDate(transactions);
  const sortedDates = Object.keys(grouped).sort().reverse();

  return (
    <div className="space-y-8">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-zinc-900 dark:text-zinc-50">Transactions</h1>
          <p className="mt-1 text-zinc-500 dark:text-zinc-400">Track your income and expenses.</p>
        </div>
        {!showForm && (
          <Button onClick={() => setShowForm(true)}>
            <Plus className="mr-1.5 h-4 w-4" />
            Add transaction
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
            onClick={() => {
              createMutation.reset();
              updateMutation.reset();
              deleteMutation.reset();
            }}
            className="rounded-lg p-1 text-red-500 transition-colors hover:bg-red-100 dark:text-red-400 dark:hover:bg-red-900/50"
            aria-label="Dismiss"
          >
            <X className="h-3.5 w-3.5" />
          </button>
        </div>
      )}

      {showForm && (
        <Card>
          <h2 className="mb-4 text-lg font-semibold text-zinc-900 dark:text-zinc-50">
            {editingId ? "Edit transaction" : "New transaction"}
          </h2>
          <form onSubmit={handleSubmit} className="flex flex-col gap-4">
            <div className="grid gap-4 sm:grid-cols-2">
              <div className="space-y-1.5">
                <label
                  htmlFor="type"
                  className="block text-sm font-medium text-zinc-700 dark:text-zinc-300"
                >
                  Type
                </label>
                <select
                  id="type"
                  value={txnType}
                  onChange={(e) => {
                    setTxnType(e.target.value);
                    if (e.target.value !== "transfer") {
                      setToAccountId("");
                    }
                  }}
                  required
                  className="block w-full rounded-xl border border-zinc-200 bg-white px-3 py-2 text-sm text-zinc-900 focus:border-zinc-400 focus:outline-none dark:border-zinc-700 dark:bg-zinc-800 dark:text-zinc-100 dark:focus:border-zinc-500"
                >
                  {TXN_TYPES.map((t) => (
                    <option key={t} value={t}>
                      {t.charAt(0).toUpperCase() + t.slice(1)}
                    </option>
                  ))}
                </select>
              </div>
              <Input
                label="Date"
                id="date"
                type="date"
                value={date}
                onChange={(e) => setDate(e.target.value)}
                required
              />
              <div className="space-y-1.5">
                <label
                  htmlFor="account"
                  className="block text-sm font-medium text-zinc-700 dark:text-zinc-300"
                >
                  {txnType === "transfer" ? "From account" : "Account"}
                </label>
                <select
                  id="account"
                  value={accountId}
                  onChange={(e) => setAccountId(e.target.value ? Number(e.target.value) : "")}
                  required
                  className="block w-full rounded-xl border border-zinc-200 bg-white px-3 py-2 text-sm text-zinc-900 focus:border-zinc-400 focus:outline-none dark:border-zinc-700 dark:bg-zinc-800 dark:text-zinc-100 dark:focus:border-zinc-500"
                >
                  <option value="">Select account</option>
                  {accounts.map((acc) => (
                    <option key={acc.id} value={acc.id}>
                      {acc.name}
                    </option>
                  ))}
                </select>
              </div>
              {txnType === "transfer" ? (
                <div className="space-y-1.5">
                  <label
                    htmlFor="toAccount"
                    className="block text-sm font-medium text-zinc-700 dark:text-zinc-300"
                  >
                    To account
                  </label>
                  <select
                    id="toAccount"
                    value={toAccountId}
                    onChange={(e) => setToAccountId(e.target.value ? Number(e.target.value) : "")}
                    required
                    className="block w-full rounded-xl border border-zinc-200 bg-white px-3 py-2 text-sm text-zinc-900 focus:border-zinc-400 focus:outline-none dark:border-zinc-700 dark:bg-zinc-800 dark:text-zinc-100 dark:focus:border-zinc-500"
                  >
                    <option value="">Select account</option>
                    {accounts
                      .filter((acc) => !accountId || acc.id !== Number(accountId))
                      .map((acc) => (
                        <option key={acc.id} value={acc.id}>
                          {acc.name}
                        </option>
                      ))}
                  </select>
                </div>
              ) : (
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
                    className="block w-full rounded-xl border border-zinc-200 bg-white px-3 py-2 text-sm text-zinc-900 focus:border-zinc-400 focus:outline-none dark:border-zinc-700 dark:bg-zinc-800 dark:text-zinc-100 dark:focus:border-zinc-500"
                  >
                    <option value="">No category</option>
                    {categories
                      .filter((c) => c.type === txnType)
                      .map((cat) => (
                        <option key={cat.id} value={cat.id}>
                          {cat.name}
                        </option>
                      ))}
                  </select>
                </div>
              )}
              <Input
                label="Amount"
                id="amount"
                type="number"
                step="0.01"
                min="0"
                value={amount}
                onChange={(e) => setAmount(e.target.value)}
                required
                placeholder="0.00"
              />
              <Input
                label="Description"
                id="description"
                value={description}
                onChange={(e) => setDescription(e.target.value)}
                required
                placeholder="e.g. Grocery shopping"
              />
            </div>
            <div className="flex justify-end gap-3">
              <Button type="button" variant="ghost" onClick={resetForm}>
                Cancel
              </Button>
              <Button
                type="submit"
                disabled={isPending || !accountId || !amount.trim() || !description.trim()}
              >
                {isPending ? "Saving..." : editingId ? "Save" : "Create"}
              </Button>
            </div>
          </form>
        </Card>
      )}

      {isLoading ? (
        <Card>
          <p className="text-center text-sm text-zinc-400 dark:text-zinc-500">
            Loading transactions...
          </p>
        </Card>
      ) : transactions.length === 0 ? (
        <Card>
          <div className="flex flex-col items-center gap-3 py-12 text-center">
            <ArrowLeftRight className="h-10 w-10 text-zinc-300 dark:text-zinc-600" />
            <p className="text-sm text-zinc-400 dark:text-zinc-500">
              No transactions yet. Create your first transaction to get started.
            </p>
          </div>
        </Card>
      ) : (
        <div className="space-y-4">
          <PaginationBar total={total} page={page} totalPages={totalPages} onPageChange={setPage} />
          <div className="space-y-6">
            {sortedDates.map((dateStr) => (
              <div key={dateStr}>
                <h2 className="mb-3 text-sm font-semibold uppercase tracking-wider text-zinc-500 dark:text-zinc-400">
                  {formatDate(dateStr)}
                </h2>
                <div className="space-y-1">
                  {grouped[dateStr].map((tx) => (
                    <div
                      key={tx.id}
                      className="group flex items-center justify-between rounded-xl border border-zinc-200 bg-white/80 px-4 py-3 shadow-sm backdrop-blur-xl dark:border-zinc-800 dark:bg-zinc-900/80"
                    >
                      <div className="flex min-w-0 items-center gap-3">
                        <span className="shrink-0 text-zinc-400 dark:text-zinc-500">
                          {tx.type === "income" ? (
                            <ArrowDownRight className="h-4 w-4 text-emerald-500" />
                          ) : tx.type === "expense" ? (
                            <ArrowUpRight className="h-4 w-4 text-red-500" />
                          ) : (
                            <ArrowRightLeft className="h-4 w-4 text-blue-500" />
                          )}
                        </span>
                        <div className="min-w-0">
                          <div className="flex items-center gap-2">
                            <span className="truncate text-sm font-medium text-zinc-900 dark:text-zinc-50">
                              {tx.description}
                            </span>
                            <Badge
                              variant={
                                tx.type === "income"
                                  ? "success"
                                  : tx.type === "expense"
                                    ? "danger"
                                    : "info"
                              }
                            >
                              {tx.type}
                            </Badge>
                          </div>
                          <div className="mt-0.5 flex items-center gap-2 text-xs text-zinc-400 dark:text-zinc-500">
                            {tx.category && (
                              <span className="flex items-center gap-1">
                                {tx.category.color && (
                                  <span
                                    className="h-2 w-2 rounded-full"
                                    style={{
                                      backgroundColor: tx.category.color,
                                    }}
                                  />
                                )}
                                {tx.category.name}
                              </span>
                            )}
                            {tx.category && tx.account && <span>&middot;</span>}
                            {tx.account && <span>{tx.account.name}</span>}
                            {tx.to_account && (
                              <>
                                <ArrowRightLeft className="h-3 w-3" />
                                <span>{tx.to_account.name}</span>
                              </>
                            )}
                          </div>
                        </div>
                      </div>
                      <div className="flex shrink-0 items-center gap-2">
                        <span
                          className={`text-sm font-semibold tabular-nums ${
                            tx.type === "income"
                              ? "text-emerald-600 dark:text-emerald-400"
                              : tx.type === "expense"
                                ? "text-red-600 dark:text-red-400"
                                : "text-zinc-900 dark:text-zinc-50"
                          }`}
                        >
                          {tx.type === "income" ? "+" : tx.type === "expense" ? "-" : ""}
                          {formatCents(tx.amount)}
                        </span>
                        <button
                          onClick={() => populateForm(tx)}
                          className="rounded-lg p-1.5 text-zinc-400 opacity-0 transition-all hover:bg-zinc-100 hover:text-zinc-600 group-hover:opacity-100 dark:hover:bg-zinc-800 dark:hover:text-zinc-300"
                          aria-label="Edit"
                        >
                          <Pencil className="h-3.5 w-3.5" />
                        </button>
                        <button
                          onClick={() => {
                            if (window.confirm(`Delete transaction "${tx.description}"?`)) {
                              deleteMutation.mutate(tx.id);
                            }
                          }}
                          className="rounded-lg p-1.5 text-zinc-400 opacity-0 transition-all hover:bg-red-100 hover:text-red-600 group-hover:opacity-100 dark:hover:bg-red-900/30 dark:hover:text-red-400"
                          aria-label="Delete"
                        >
                          <Trash2 className="h-3.5 w-3.5" />
                        </button>
                      </div>
                    </div>
                  ))}
                </div>
              </div>
            ))}
          </div>
          <PaginationBar total={total} page={page} totalPages={totalPages} onPageChange={setPage} />
        </div>
      )}
    </div>
  );
}

function PaginationBar({
  total,
  page,
  totalPages,
  onPageChange,
}: {
  total: number;
  page: number;
  totalPages: number;
  onPageChange: (p: number) => void;
}) {
  if (total === 0) return null;
  return (
    <div className="flex items-center justify-between">
      <p className="text-xs text-zinc-400 dark:text-zinc-500">
        {total} transaction{total !== 1 ? "s" : ""}
      </p>
      <div className="flex items-center gap-2">
        <Button
          variant="ghost"
          size="sm"
          disabled={page <= 1}
          onClick={() => onPageChange(Math.max(1, page - 1))}
        >
          Previous
        </Button>
        {totalPages > 1 &&
          Array.from({ length: Math.min(totalPages, 5) }, (_, i) => {
            const start = Math.max(1, page - 2);
            const p = start + i;
            if (p > totalPages) return null;
            return (
              <button
                key={p}
                onClick={() => onPageChange(p)}
                className={`inline-flex h-8 w-8 items-center justify-center rounded-xl text-sm font-medium transition-colors ${
                  p === page
                    ? "bg-zinc-900 text-white dark:bg-zinc-100 dark:text-zinc-900"
                    : "text-zinc-600 hover:bg-zinc-100 dark:text-zinc-400 dark:hover:bg-zinc-800"
                }`}
              >
                {p}
              </button>
            );
          })}
        <Button
          variant="ghost"
          size="sm"
          disabled={page >= totalPages}
          onClick={() => onPageChange(page + 1)}
        >
          Next
        </Button>
      </div>
    </div>
  );
}

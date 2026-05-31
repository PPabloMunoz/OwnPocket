import { useState, type FormEvent } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { Plus, Wallet, Pencil, Trash2, X, AlertCircle } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Card } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { api } from "@/lib/api";
import { queryKeys } from "@/lib/query-keys";
import { formatCents } from "@/lib/utils";
import type { Account } from "@/types/account";

export function meta() {
  return [{ title: "OwnPocket - Accounts" }];
}

const ACCOUNT_TYPES = ["checking", "savings", "credit_card", "cash", "investment", "loan"] as const;

const typeBadgeVariant: Record<string, "info" | "success" | "warning" | "danger" | "default"> = {
  checking: "info",
  savings: "success",
  credit_card: "warning",
  cash: "default",
  investment: "info",
  loan: "danger",
};

export default function AccountsPage() {
  const queryClient = useQueryClient();
  const [showForm, setShowForm] = useState(false);
  const [editingId, setEditingId] = useState<number | null>(null);
  const [name, setName] = useState("");
  const [type, setType] = useState<string>("checking");
  const [description, setDescription] = useState("");
  const [initialBalance, setInitialBalance] = useState("");

  const { data: accounts = [], isLoading } = useQuery({
    queryKey: queryKeys.accounts.all,
    queryFn: () => api.get<Account[]>("/accounts"),
  });

  const createMutation = useMutation({
    mutationFn: (body: { name: string; type: string; description?: string; balance?: number }) =>
      api.post<Account>("/accounts", body),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.accounts.all });
      resetForm();
    },
  });

  const updateMutation = useMutation({
    mutationFn: ({
      id,
      ...body
    }: {
      id: number;
      name: string;
      type: string;
      description?: string;
      balance?: number;
    }) => api.put<Account>(`/accounts/${id}`, body),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.accounts.all });
      resetForm();
    },
  });

  const deleteMutation = useMutation({
    mutationFn: (id: number) => api.delete(`/accounts/${id}`),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.accounts.all });
    },
  });

  function resetForm() {
    setShowForm(false);
    setEditingId(null);
    setName("");
    setType("checking");
    setDescription("");
    setInitialBalance("");
  }

  function populateForm(account: Account) {
    setEditingId(account.id);
    setName(account.name);
    setType(account.type);
    setDescription(account.description ?? "");
    setInitialBalance(account.balance ? (account.balance / 100).toFixed(2) : "");
    setShowForm(true);
  }

  const handleSubmit = (e: FormEvent) => {
    e.preventDefault();
    if (!name.trim()) return;
    const balance = initialBalance ? parseFloat(initialBalance) : undefined;
    const payload: {
      name: string;
      type: string;
      description?: string;
      balance?: number;
    } = { name: name.trim(), type };
    if (description.trim()) payload.description = description.trim();
    if (balance !== undefined) payload.balance = balance;
    if (editingId) {
      updateMutation.mutate({ id: editingId, ...payload });
    } else {
      createMutation.mutate(payload);
    }
  };

  const isPending = createMutation.isPending || updateMutation.isPending;
  const mutationError = createMutation.error ?? updateMutation.error ?? deleteMutation.error;

  return (
    <div className="space-y-8">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-zinc-900 dark:text-zinc-50">Accounts</h1>
          <p className="mt-1 text-zinc-500 dark:text-zinc-400">
            Manage your bank accounts and wallets.
          </p>
        </div>
        {!showForm && (
          <Button onClick={() => setShowForm(true)}>
            <Plus className="mr-1.5 h-4 w-4" />
            Add account
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
            {editingId ? "Edit account" : "New account"}
          </h2>
          <form onSubmit={handleSubmit} className="flex flex-col gap-4">
            <div className="grid gap-4 sm:grid-cols-2">
              <Input
                label="Name"
                id="name"
                value={name}
                onChange={(e) => setName(e.target.value)}
                required
                placeholder="e.g. Main Checking"
              />
              <div className="space-y-1.5">
                <label
                  htmlFor="type"
                  className="block text-sm font-medium text-zinc-700 dark:text-zinc-300"
                >
                  Type
                </label>
                <select
                  id="type"
                  value={type}
                  onChange={(e) => setType(e.target.value)}
                  required
                  className="block w-full rounded-xl border border-zinc-200 bg-white px-3 py-2 text-sm text-zinc-900 focus:border-zinc-400 focus:outline-none dark:border-zinc-700 dark:bg-zinc-800 dark:text-zinc-100 dark:focus:border-zinc-500"
                >
                  {ACCOUNT_TYPES.map((t) => (
                    <option key={t} value={t}>
                      {t.replace("_", " ").replace(/\b\w/g, (c) => c.toUpperCase())}
                    </option>
                  ))}
                </select>
              </div>
              <Input
                label="Description"
                id="description"
                value={description}
                onChange={(e) => setDescription(e.target.value)}
                placeholder="Optional note"
              />
              <Input
                label={editingId ? "Balance" : "Initial balance"}
                id="balance"
                type="number"
                step="0.01"
                value={initialBalance}
                onChange={(e) => setInitialBalance(e.target.value)}
                placeholder="0.00"
              />
            </div>
            <div className="flex justify-end gap-3">
              <Button type="button" variant="ghost" onClick={resetForm}>
                Cancel
              </Button>
              <Button type="submit" disabled={isPending || !name.trim()}>
                {isPending ? "Saving..." : editingId ? "Save" : "Create"}
              </Button>
            </div>
          </form>
        </Card>
      )}

      {isLoading ? (
        <Card>
          <p className="text-center text-sm text-zinc-400 dark:text-zinc-500">
            Loading accounts...
          </p>
        </Card>
      ) : accounts.length === 0 ? (
        <Card>
          <div className="flex flex-col items-center gap-3 py-12 text-center">
            <Wallet className="h-10 w-10 text-zinc-300 dark:text-zinc-600" />
            <p className="text-sm text-zinc-400 dark:text-zinc-500">
              No accounts yet. Add your first account to start tracking.
            </p>
          </div>
        </Card>
      ) : (
        <div className="space-y-1">
          {accounts.map((account) => (
            <div
              key={account.id}
              className="group flex items-center justify-between rounded-xl border border-zinc-200 bg-white/80 px-4 py-3 shadow-sm backdrop-blur-xl dark:border-zinc-800 dark:bg-zinc-900/80"
            >
              <div className="flex items-center gap-3">
                <div className="flex flex-col">
                  <div className="flex items-center gap-2">
                    <span className="text-sm font-medium text-zinc-900 dark:text-zinc-50">
                      {account.name}
                    </span>
                    <Badge variant={typeBadgeVariant[account.type]}>
                      {account.type.replace("_", " ")}
                    </Badge>
                  </div>
                  {account.description && (
                    <p className="mt-0.5 text-xs text-zinc-400 dark:text-zinc-500">
                      {account.description}
                    </p>
                  )}
                </div>
              </div>
              <div className="flex items-center gap-2">
                <span className="text-sm font-semibold tabular-nums text-zinc-900 dark:text-zinc-50">
                  {formatCents(account.balance)}
                </span>
                <button
                  onClick={() => populateForm(account)}
                  className="rounded-lg p-1.5 text-zinc-400 opacity-0 transition-all hover:bg-zinc-100 hover:text-zinc-600 group-hover:opacity-100 dark:hover:bg-zinc-800 dark:hover:text-zinc-300"
                  aria-label="Edit"
                >
                  <Pencil className="h-3.5 w-3.5" />
                </button>
                <button
                  onClick={() => {
                    if (window.confirm(`Delete account "${account.name}"?`)) {
                      deleteMutation.mutate(account.id);
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
      )}
    </div>
  );
}

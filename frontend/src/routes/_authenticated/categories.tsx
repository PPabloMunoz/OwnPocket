import { createFileRoute } from "@tanstack/react-router";
import { useState, type FormEvent } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { Plus, Tags, FolderOpen } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Badge } from "@/components/ui/badge";
import { Dialog } from "@/components/ui/dialog";
import { api } from "@/lib/api";
import { queryKeys } from "@/lib/query-keys";
import type { Category, CreateCategoryRequest } from "@/types/category";

export const Route = createFileRoute("/_authenticated/categories")({
  component: CategoriesPage,
});

function CategoriesPage() {
  const queryClient = useQueryClient();
  const [showForm, setShowForm] = useState(false);
  const [name, setName] = useState("");
  const [type, setType] = useState<"income" | "expense">("expense");
  const [parentId, setParentId] = useState<number | "">("");
  const [color, setColor] = useState("");

  const { data: categories = [], isLoading } = useQuery({
    queryKey: queryKeys.categories.all,
    queryFn: () => api.get<Category[]>("/categories"),
  });

  const createMutation = useMutation({
    mutationFn: (body: CreateCategoryRequest) =>
      api.post<Category>("/categories", body),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.categories.all });
      setName("");
      setType("expense");
      setParentId("");
      setColor("");
      setShowForm(false);
    },
  });

  const handleSubmit = (e: FormEvent) => {
    e.preventDefault();
    if (!name.trim()) return;
    createMutation.mutate({
      name: name.trim(),
      type,
      ...(parentId !== "" ? { parent_id: parentId } : {}),
      ...(color.trim() ? { color: color.trim() } : {}),
    });
  };

  const incomeCategories = categories.filter((c) => c.type === "income");
  const expenseCategories = categories.filter((c) => c.type === "expense");
  const parentOptions = categories.filter((c) => c.type === type && !c.parent_id);

  return (
    <div className="space-y-8">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-zinc-900 dark:text-zinc-50">
            Categories
          </h1>
          <p className="mt-1 text-zinc-500 dark:text-zinc-400">
            Organize your transactions by category.
          </p>
        </div>
        <Button onClick={() => setShowForm(true)}>
          <Plus className="mr-1.5 h-4 w-4" />
          Add category
        </Button>
      </div>

      {isLoading ? (
        <div className="rounded-2xl border border-zinc-200 bg-white/80 p-6 shadow-sm backdrop-blur-xl dark:border-zinc-800 dark:bg-zinc-900/80">
          <p className="text-center text-sm text-zinc-400 dark:text-zinc-500">
            Loading categories...
          </p>
        </div>
      ) : categories.length === 0 ? (
        <div className="rounded-2xl border border-zinc-200 bg-white/80 p-6 shadow-sm backdrop-blur-xl dark:border-zinc-800 dark:bg-zinc-900/80">
          <div className="flex flex-col items-center gap-3 py-12 text-center">
            <Tags className="h-10 w-10 text-zinc-300 dark:text-zinc-600" />
            <p className="text-sm text-zinc-400 dark:text-zinc-500">
              No categories yet. Add categories to organize your transactions.
            </p>
          </div>
        </div>
      ) : (
        <div className="space-y-6">
          {incomeCategories.length > 0 && (
            <CategoryGroup title="Income" categories={incomeCategories} />
          )}
          {expenseCategories.length > 0 && (
            <CategoryGroup title="Expenses" categories={expenseCategories} />
          )}
        </div>
      )}

      <Dialog open={showForm} onClose={() => setShowForm(false)} title="New Category">
        <form onSubmit={handleSubmit} className="flex flex-col gap-4">
          <Input
            label="Name"
            id="name"
            value={name}
            onChange={(e) => setName(e.target.value)}
            required
            placeholder="e.g. Groceries"
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
              onChange={(e) => {
                setType(e.target.value as "income" | "expense");
                setParentId("");
              }}
              className="block w-full rounded-xl border border-zinc-200 bg-white px-3 py-2 text-sm text-zinc-900 focus:border-zinc-400 focus:outline-none dark:border-zinc-700 dark:bg-zinc-800 dark:text-zinc-100 dark:focus:border-zinc-500"
            >
              <option value="expense">Expense</option>
              <option value="income">Income</option>
            </select>
          </div>
          {parentOptions.length > 0 && (
            <div className="space-y-1.5">
              <label
                htmlFor="parent"
                className="block text-sm font-medium text-zinc-700 dark:text-zinc-300"
              >
                Parent category
              </label>
              <select
                id="parent"
                value={parentId}
                onChange={(e) => setParentId(e.target.value ? Number(e.target.value) : "")}
                className="block w-full rounded-xl border border-zinc-200 bg-white px-3 py-2 text-sm text-zinc-900 focus:border-zinc-400 focus:outline-none dark:border-zinc-700 dark:bg-zinc-800 dark:text-zinc-100 dark:focus:border-zinc-500"
              >
                <option value="">None</option>
                {parentOptions.map((cat) => (
                  <option key={cat.id} value={cat.id}>{cat.name}</option>
                ))}
              </select>
            </div>
          )}
          <Input
            label="Color"
            id="color"
            value={color}
            onChange={(e) => setColor(e.target.value)}
            placeholder="e.g. #10b981"
          />
          {createMutation.error && (
            <p className="text-sm text-red-600 dark:text-red-400">
              {createMutation.error instanceof Error
                ? createMutation.error.message
                : "Failed to create category"}
            </p>
          )}
          <div className="flex justify-end gap-3 pt-2">
            <Button type="button" variant="ghost" onClick={() => setShowForm(false)}>
              Cancel
            </Button>
            <Button type="submit" disabled={createMutation.isPending || !name.trim()}>
              {createMutation.isPending ? "Creating..." : "Create"}
            </Button>
          </div>
        </form>
      </Dialog>
    </div>
  );
}

function CategoryGroup({ title, categories }: { title: string; categories: Category[] }) {
  return (
    <div>
      <h2 className="mb-3 text-sm font-semibold uppercase tracking-wider text-zinc-500 dark:text-zinc-400">
        {title}
      </h2>
      <div className="space-y-1">
        {categories.map((cat) => (
          <div
            key={cat.id}
            className="flex items-center justify-between rounded-xl border border-zinc-200 bg-white/80 px-4 py-3 shadow-sm backdrop-blur-xl dark:border-zinc-800 dark:bg-zinc-900/80"
          >
            <div className="flex items-center gap-3">
              {cat.color ? (
                <span
                  className="h-3 w-3 shrink-0 rounded-full"
                  style={{ backgroundColor: cat.color }}
                />
              ) : (
                <FolderOpen className="h-3.5 w-3.5 text-zinc-400 dark:text-zinc-500" />
              )}
              <div>
                <span className="text-sm font-medium text-zinc-900 dark:text-zinc-50">
                  {cat.name}
                </span>
                {cat.parent && (
                  <span className="ml-2 text-xs text-zinc-400 dark:text-zinc-500">
                    {cat.parent.name}
                  </span>
                )}
              </div>
            </div>
            <Badge variant={cat.type === "income" ? "success" : "info"}>
              {cat.type}
            </Badge>
          </div>
        ))}
      </div>
    </div>
  );
}

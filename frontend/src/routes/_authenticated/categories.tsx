import { createFileRoute } from "@tanstack/react-router";
import { useState, type FormEvent } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { Plus, Tags, Pencil, Trash2, X, AlertCircle } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Badge } from "@/components/ui/badge";
import { api } from "@/lib/api";
import { queryKeys } from "@/lib/query-keys";
import type { Category, CreateCategoryRequest } from "@/types/category";

export const Route = createFileRoute("/_authenticated/categories")({
  component: CategoriesPage,
});

function buildTree(categories: Category[], parentId: number | null = null): Category[] {
  return categories
    .filter((c) => c.parent_id === parentId)
    .sort((a, b) => a.name.localeCompare(b.name));
}

function flattenTree(
  categories: Category[],
  roots: Category[],
  depth = 0,
): { cat: Category; depth: number }[] {
  const result: { cat: Category; depth: number }[] = [];
  for (const root of roots) {
    result.push({ cat: root, depth });
    const children = buildTree(categories, root.id);
    result.push(...flattenTree(categories, children, depth + 1));
  }
  return result;
}

function collectDescendantIds(categories: Category[], id: number): Set<number> {
  const ids = new Set<number>();
  const stack = [id];
  while (stack.length > 0) {
    const current = stack.pop()!;
    for (const cat of categories) {
      if (cat.parent_id === current && !ids.has(cat.id)) {
        ids.add(cat.id);
        stack.push(cat.id);
      }
    }
  }
  return ids;
}

function CategoriesPage() {
  const queryClient = useQueryClient();
  const [showForm, setShowForm] = useState(false);
  const [editingId, setEditingId] = useState<number | null>(null);
  const [name, setName] = useState("");
  const [type, setType] = useState<"income" | "expense">("expense");
  const [parentId, setParentId] = useState<number | "">("");
  const [color, setColor] = useState("");

  const { data: categories = [], isLoading } = useQuery({
    queryKey: queryKeys.categories.all,
    queryFn: () => api.get<Category[]>("/categories"),
  });

  const createMutation = useMutation({
    mutationFn: (body: CreateCategoryRequest) => api.post<Category>("/categories", body),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.categories.all });
      resetForm();
    },
  });

  const updateMutation = useMutation({
    mutationFn: ({ id, ...body }: CreateCategoryRequest & { id: number }) =>
      api.put<Category>(`/categories/${id}`, body),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.categories.all });
      resetForm();
    },
  });

  const deleteMutation = useMutation({
    mutationFn: (id: number) => api.delete(`/categories/${id}`),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.categories.all });
    },
  });

  function resetForm() {
    setShowForm(false);
    setEditingId(null);
    setName("");
    setType("expense");
    setParentId("");
    setColor("");
  }

  function populateForm(cat: Category) {
    setEditingId(cat.id);
    setName(cat.name);
    setType(cat.type);
    setParentId(cat.parent_id ?? "");
    setColor(cat.color ?? "");
    setShowForm(true);
  }

  const handleSubmit = (e: FormEvent) => {
    e.preventDefault();
    if (!name.trim()) return;
    const body: CreateCategoryRequest = {
      name: name.trim(),
      type,
      parent_id: parentId !== "" ? parentId : null,
      ...(color.trim() ? { color: color.trim() } : {}),
    };
    if (editingId) {
      updateMutation.mutate({ id: editingId, ...body });
    } else {
      createMutation.mutate(body);
    }
  };

  const isPending = createMutation.isPending || updateMutation.isPending;
  const mutationError = createMutation.error ?? updateMutation.error ?? deleteMutation.error;
  const incomeRoots = buildTree(categories, null).filter((c) => c.type === "income");
  const expenseRoots = buildTree(categories, null).filter((c) => c.type === "expense");
  const incomeFlat = flattenTree(categories, incomeRoots);
  const expenseFlat = flattenTree(categories, expenseRoots);

  return (
    <div className="space-y-8">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-zinc-900 dark:text-zinc-50">Categories</h1>
          <p className="mt-1 text-zinc-500 dark:text-zinc-400">
            Organize your transactions by category.
          </p>
        </div>
        {!showForm && (
          <Button onClick={() => setShowForm(true)}>
            <Plus className="mr-1.5 h-4 w-4" />
            Add category
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
        <div className="rounded-2xl border border-zinc-200 bg-white/80 p-6 shadow-sm backdrop-blur-xl dark:border-zinc-800 dark:bg-zinc-900/80">
          <h2 className="mb-4 text-lg font-semibold text-zinc-900 dark:text-zinc-50">
            {editingId ? "Edit category" : "New category"}
          </h2>
          <form onSubmit={handleSubmit} className="flex flex-col gap-4">
            <div className="grid gap-4 sm:grid-cols-2">
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
              <ParentSelect
                categories={categories}
                type={type}
                value={parentId}
                onChange={setParentId}
                editingId={editingId}
              />
              <Input
                label="Color"
                id="color"
                value={color}
                onChange={(e) => setColor(e.target.value)}
                placeholder="e.g. #10b981"
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
        </div>
      )}

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
          {incomeFlat.length > 0 && (
            <CategoryTree
              title="Income"
              flat={incomeFlat}
              onEdit={populateForm}
              onDelete={(id) => {
                if (window.confirm("Delete this category?")) {
                  deleteMutation.mutate(id);
                }
              }}
            />
          )}
          {expenseFlat.length > 0 && (
            <CategoryTree
              title="Expenses"
              flat={expenseFlat}
              onEdit={populateForm}
              onDelete={(id) => {
                if (window.confirm("Delete this category?")) {
                  deleteMutation.mutate(id);
                }
              }}
            />
          )}
        </div>
      )}
    </div>
  );
}

function ParentSelect({
  categories,
  type,
  value,
  onChange,
  editingId,
}: {
  categories: Category[];
  type: "income" | "expense";
  value: number | "";
  onChange: (v: number | "") => void;
  editingId?: number | null;
}) {
  const excludeIds = new Set<number>();
  if (editingId) {
    excludeIds.add(editingId);
    for (const id of collectDescendantIds(categories, editingId)) {
      excludeIds.add(id);
    }
  }
  const available = categories.filter(
    (c) => c.type === type && c.parent_id === null && !excludeIds.has(c.id),
  );
  if (available.length === 0) return null;
  const flat = flattenTree(categories, available).filter(({ cat }) => !excludeIds.has(cat.id));

  return (
    <div className="space-y-1.5">
      <label
        htmlFor="parent"
        className="block text-sm font-medium text-zinc-700 dark:text-zinc-300"
      >
        Parent category
      </label>
      <select
        id="parent"
        value={value}
        onChange={(e) => onChange(e.target.value ? Number(e.target.value) : "")}
        className="block w-full rounded-xl border border-zinc-200 bg-white px-3 py-2 text-sm text-zinc-900 focus:border-zinc-400 focus:outline-none dark:border-zinc-700 dark:bg-zinc-800 dark:text-zinc-100 dark:focus:border-zinc-500"
      >
        <option value="">None</option>
        {flat.map(({ cat, depth }) => (
          <option key={cat.id} value={cat.id}>
            {"\u00A0".repeat(depth * 4)}
            {depth > 0 ? "\u2514\u2500 " : ""}
            {cat.name}
          </option>
        ))}
      </select>
    </div>
  );
}

function CategoryTree({
  title,
  flat,
  onEdit,
  onDelete,
}: {
  title: string;
  flat: { cat: Category; depth: number }[];
  onEdit: (cat: Category) => void;
  onDelete: (id: number) => void;
}) {
  return (
    <div>
      <h2 className="mb-3 text-sm font-semibold uppercase tracking-wider text-zinc-500 dark:text-zinc-400">
        {title}
      </h2>
      <div className="space-y-1">
        {flat.map(({ cat, depth }) => (
          <div
            key={cat.id}
            className="group flex items-center justify-between rounded-xl border border-zinc-200 bg-white/80 px-4 py-3 shadow-sm backdrop-blur-xl dark:border-zinc-800 dark:bg-zinc-900/80"
          >
            <div className="flex items-center gap-2" style={{ paddingLeft: `${depth * 1.5}rem` }}>
              <span
                className="h-3 w-3 shrink-0 rounded-full"
                style={{ backgroundColor: cat.color ?? "oklch(0.5 0 0)" }}
              />
              <span className="text-sm font-medium text-zinc-900 dark:text-zinc-50">
                {cat.name}
              </span>
            </div>
            <div className="flex items-center gap-2">
              <button
                onClick={() => onEdit(cat)}
                className="rounded-lg p-1.5 text-zinc-400 opacity-0 transition-all hover:bg-zinc-100 hover:text-zinc-600 group-hover:opacity-100 dark:hover:bg-zinc-800 dark:hover:text-zinc-300"
                aria-label="Edit"
              >
                <Pencil className="h-3.5 w-3.5" />
              </button>
              <button
                onClick={() => onDelete(cat.id)}
                className="rounded-lg p-1.5 text-zinc-400 opacity-0 transition-all hover:bg-red-100 hover:text-red-600 group-hover:opacity-100 dark:hover:bg-red-900/30 dark:hover:text-red-400"
                aria-label="Delete"
              >
                <Trash2 className="h-3.5 w-3.5" />
              </button>
              <Badge variant={cat.type === "income" ? "success" : "info"}>{cat.type}</Badge>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}

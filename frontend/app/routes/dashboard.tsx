import { useQuery } from "@tanstack/react-query";
import {
  ArrowDownRight,
  ArrowUpRight,
  ArrowRightLeft,
  TrendingUp,
  TrendingDown,
  PiggyBank,
} from "lucide-react";
import { api } from "@/lib/api";
import { queryKeys } from "@/lib/query-keys";
import { formatCents } from "@/lib/utils";
import type { DashboardSummary } from "@/types/dashboard";

export function meta() {
  return [{ title: "OwnPocket - Dashboard" }];
}

function pct(a: number, b: number): number {
  if (b === 0) return 0;
  return Math.min(Math.round((a / b) * 100), 100);
}

export default function DashboardPage() {
  const { data: summary, isLoading } = useQuery({
    queryKey: queryKeys.dashboard.summary,
    queryFn: () => api.get<DashboardSummary>("/dashboard/summary"),
  });

  const balance = summary?.total_balance ?? 0;
  const income = summary?.monthly_income ?? 0;
  const expenses = summary?.monthly_expenses ?? 0;
  const netMonth = income - expenses;
  const txs = summary?.recent_transactions ?? [];
  const budgets = summary?.budgets ?? [];

  return (
    <div className="space-y-8">
      {/* Hero */}
      <section style={{ animationDelay: "0ms" }} className="animate-[fadeIn_0.6s_ease_both]">
        <p className="mb-1 text-xs font-medium uppercase tracking-[0.15em] text-zinc-400 dark:text-zinc-500">
          Total balance
        </p>
        <p className="text-4xl font-bold tracking-tight text-zinc-900 dark:text-zinc-50 sm:text-5xl">
          {isLoading ? (
            <span className="text-zinc-300 dark:text-zinc-700">—</span>
          ) : (
            formatCents(balance)
          )}
        </p>
      </section>

      {/* KPI strip */}
      <section
        style={{ animationDelay: "100ms" }}
        className="animate-[fadeIn_0.6s_ease_both] grid gap-3 sm:grid-cols-3"
      >
        {(
          [
            {
              label: "Monthly income",
              value: income,
              icon: TrendingUp,
              color: "text-emerald-500",
              bg: "bg-emerald-500/10",
            },
            {
              label: "Monthly expenses",
              value: expenses,
              icon: TrendingDown,
              color: "text-rose-500",
              bg: "bg-rose-500/10",
            },
            {
              label: "Net this month",
              value: netMonth,
              icon: PiggyBank,
              color: netMonth >= 0 ? "text-emerald-500" : "text-rose-500",
              bg: "bg-zinc-500/10",
            },
          ] as const
        ).map((kpi) => (
          <div
            key={kpi.label}
            className="rounded-2xl border border-zinc-200 bg-white/80 p-4 shadow-sm backdrop-blur-xl dark:border-zinc-800 dark:bg-zinc-900/80"
          >
            <div className="flex items-center gap-3">
              <div
                className={`flex h-9 w-9 shrink-0 items-center justify-center rounded-xl ${kpi.bg}`}
              >
                <kpi.icon className={`h-4 w-4 ${kpi.color}`} />
              </div>
              <div className="min-w-0">
                <p className="truncate text-xs text-zinc-400 dark:text-zinc-500">{kpi.label}</p>
                <p className={`text-base font-bold tabular-nums ${kpi.value >= 0 ? "" : ""}`}>
                  {isLoading ? (
                    <span className="text-zinc-300 dark:text-zinc-700">—</span>
                  ) : (
                    formatCents(kpi.value)
                  )}
                </p>
              </div>
            </div>
          </div>
        ))}
      </section>

      {/* Two-column content */}
      <div className="grid gap-6 lg:grid-cols-3">
        {/* Recent transactions */}
        <section
          style={{ animationDelay: "200ms" }}
          className="animate-[fadeIn_0.6s_ease_both] lg:col-span-2"
        >
          <div className="rounded-2xl border border-zinc-200 bg-white/80 shadow-sm backdrop-blur-xl dark:border-zinc-800 dark:bg-zinc-900/80">
            <div className="border-b border-zinc-200 px-5 py-3.5 dark:border-zinc-800">
              <h2 className="text-sm font-semibold text-zinc-900 dark:text-zinc-50">
                Recent transactions
              </h2>
            </div>
            {isLoading ? (
              <div className="px-5 py-10 text-center text-sm text-zinc-400">Loading...</div>
            ) : txs.length === 0 ? (
              <div className="flex flex-col items-center gap-2 px-5 py-10 text-center">
                <ArrowRightLeft className="h-6 w-6 text-zinc-300 dark:text-zinc-600" />
                <p className="text-sm text-zinc-400 dark:text-zinc-500">No transactions yet.</p>
              </div>
            ) : (
              <div className="divide-y divide-zinc-100 dark:divide-zinc-800">
                {txs.slice(0, 7).map((tx, i) => (
                  <div
                    key={tx.id}
                    className="animate-[fadeIn_0.6s_ease_both] flex items-center justify-between px-5 py-3 transition-colors hover:bg-zinc-50 dark:hover:bg-zinc-800/50"
                    style={{ animationDelay: `${300 + i * 40}ms` }}
                  >
                    <div className="flex min-w-0 items-center gap-3">
                      <span className="shrink-0">
                        {tx.type === "income" ? (
                          <ArrowDownRight className="h-4 w-4 text-emerald-500" />
                        ) : tx.type === "expense" ? (
                          <ArrowUpRight className="h-4 w-4 text-rose-500" />
                        ) : (
                          <ArrowRightLeft className="h-4 w-4 text-blue-500" />
                        )}
                      </span>
                      <div className="min-w-0">
                        <p className="truncate text-sm font-medium text-zinc-900 dark:text-zinc-50">
                          {tx.description}
                        </p>
                        <p className="truncate text-xs text-zinc-400 dark:text-zinc-500">
                          {tx.category?.name}
                          {tx.category && tx.account && <span className="mx-1">&middot;</span>}
                          {tx.account?.name}
                        </p>
                      </div>
                    </div>
                    <span
                      className={`shrink-0 text-sm font-semibold tabular-nums ${
                        tx.type === "income"
                          ? "text-emerald-600 dark:text-emerald-400"
                          : tx.type === "expense"
                            ? "text-rose-600 dark:text-rose-400"
                            : "text-zinc-900 dark:text-zinc-50"
                      }`}
                    >
                      {tx.type === "income" ? "+" : tx.type === "expense" ? "−" : ""}
                      {formatCents(tx.amount)}
                    </span>
                  </div>
                ))}
              </div>
            )}
          </div>
        </section>

        {/* Budget progress */}
        <section style={{ animationDelay: "300ms" }} className="animate-[fadeIn_0.6s_ease_both]">
          <div className="rounded-2xl border border-zinc-200 bg-white/80 shadow-sm backdrop-blur-xl dark:border-zinc-800 dark:bg-zinc-900/80">
            <div className="border-b border-zinc-200 px-5 py-3.5 dark:border-zinc-800">
              <h2 className="text-sm font-semibold text-zinc-900 dark:text-zinc-50">
                Budget progress
              </h2>
            </div>
            {isLoading ? (
              <div className="px-5 py-10 text-center text-sm text-zinc-400">Loading...</div>
            ) : budgets.length === 0 ? (
              <div className="flex flex-col items-center gap-2 px-5 py-10 text-center">
                <PiggyBank className="h-6 w-6 text-zinc-300 dark:text-zinc-600" />
                <p className="text-sm text-zinc-400 dark:text-zinc-500">No budgets set.</p>
              </div>
            ) : (
              <div className="space-y-4 px-5 py-4">
                {budgets.slice(0, 5).map((b, i) => {
                  const used = pct(b.spent, b.amount);
                  const over = b.spent > b.amount;
                  return (
                    <div
                      key={b.id}
                      className="animate-[fadeIn_0.6s_ease_both]"
                      style={{ animationDelay: `${400 + i * 60}ms` }}
                    >
                      <div className="mb-1.5 flex items-center justify-between text-xs">
                        <span className="font-medium text-zinc-700 dark:text-zinc-300">
                          {b.category?.name ?? "Unknown"}
                        </span>
                        <span
                          className={`tabular-nums ${
                            over ? "text-rose-500" : "text-zinc-500 dark:text-zinc-400"
                          }`}
                        >
                          {formatCents(b.spent)}
                          <span className="text-zinc-300 dark:text-zinc-600">
                            {" "}
                            / {formatCents(b.amount)}
                          </span>
                        </span>
                      </div>
                      <div className="h-1.5 overflow-hidden rounded-full bg-zinc-100 dark:bg-zinc-800">
                        <div
                          className={`h-full rounded-full transition-all duration-700 ${
                            over ? "bg-rose-500" : used > 80 ? "bg-amber-500" : "bg-emerald-500"
                          }`}
                          style={{ width: `${Math.min(used, 100)}%` }}
                        />
                      </div>
                    </div>
                  );
                })}
              </div>
            )}
          </div>
        </section>
      </div>
    </div>
  );
}

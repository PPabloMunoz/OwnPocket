import { createFileRoute } from "@tanstack/react-router";
import { Wallet, TrendingUp, TrendingDown, ArrowRight } from "lucide-react";

export const Route = createFileRoute("/_authenticated/dashboard")({
  component: DashboardPage,
});

const summaryCards = [
  {
    label: "Total Balance",
    value: "$0.00",
    icon: Wallet,
    color: "text-zinc-900 dark:text-zinc-100",
    bg: "bg-zinc-900/10 dark:bg-zinc-100/10",
  },
  {
    label: "Monthly Income",
    value: "$0.00",
    icon: TrendingUp,
    color: "text-emerald-600 dark:text-emerald-400",
    bg: "bg-emerald-500/10",
  },
  {
    label: "Monthly Expenses",
    value: "$0.00",
    icon: TrendingDown,
    color: "text-red-600 dark:text-red-400",
    bg: "bg-red-500/10",
  },
];

function DashboardPage() {
  return (
    <div className="space-y-8">
      <div>
        <h1 className="text-2xl font-bold text-zinc-900 dark:text-zinc-50">
          Dashboard
        </h1>
        <p className="mt-1 text-zinc-500 dark:text-zinc-400">
          Overview of your finances.
        </p>
      </div>

      <div className="grid gap-4 sm:grid-cols-3">
        {summaryCards.map((card) => (
          <div
            key={card.label}
            className="rounded-2xl border border-zinc-200 bg-white/80 p-5 shadow-sm backdrop-blur-xl dark:border-zinc-800 dark:bg-zinc-900/80"
          >
            <div className="flex items-center gap-3">
              <div className={`flex h-10 w-10 items-center justify-center rounded-xl ${card.bg}`}>
                <card.icon className={`h-5 w-5 ${card.color}`} />
              </div>
              <div>
                <p className="text-sm text-zinc-500 dark:text-zinc-400">
                  {card.label}
                </p>
                <p className="text-xl font-bold text-zinc-900 dark:text-zinc-50">
                  {card.value}
                </p>
              </div>
            </div>
          </div>
        ))}
      </div>

      <div className="rounded-2xl border border-zinc-200 bg-white/80 p-6 shadow-sm backdrop-blur-xl dark:border-zinc-800 dark:bg-zinc-900/80">
        <div className="flex items-center justify-between">
          <h2 className="text-lg font-semibold text-zinc-900 dark:text-zinc-50">
            Recent Transactions
          </h2>
          <span className="flex items-center gap-1 text-sm text-zinc-400 dark:text-zinc-500">
            No transactions yet
          </span>
        </div>
        <div className="mt-6 flex flex-col items-center gap-2 py-8 text-center">
          <ArrowRight className="h-8 w-8 text-zinc-300 dark:text-zinc-600" />
          <p className="text-sm text-zinc-400 dark:text-zinc-500">
            Create your first transaction to get started.
          </p>
        </div>
      </div>
    </div>
  );
}

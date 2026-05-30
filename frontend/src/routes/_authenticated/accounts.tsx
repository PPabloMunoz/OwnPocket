import { createFileRoute } from "@tanstack/react-router";
import { Plus } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Wallet } from "lucide-react";

export const Route = createFileRoute("/_authenticated/accounts")({
  component: AccountsPage,
});

function AccountsPage() {
  return (
    <div className="space-y-8">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-zinc-900 dark:text-zinc-50">
            Accounts
          </h1>
          <p className="mt-1 text-zinc-500 dark:text-zinc-400">
            Manage your bank accounts and wallets.
          </p>
        </div>
        <Button>
          <Plus className="mr-1.5 h-4 w-4" />
          Add account
        </Button>
      </div>

      <div className="rounded-2xl border border-zinc-200 bg-white/80 p-6 shadow-sm backdrop-blur-xl dark:border-zinc-800 dark:bg-zinc-900/80">
        <div className="flex flex-col items-center gap-3 py-12 text-center">
          <Wallet className="h-10 w-10 text-zinc-300 dark:text-zinc-600" />
          <p className="text-sm text-zinc-400 dark:text-zinc-500">
            No accounts yet. Add your first account to start tracking.
          </p>
        </div>
      </div>
    </div>
  );
}

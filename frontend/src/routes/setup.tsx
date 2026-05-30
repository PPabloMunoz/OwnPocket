import { createFileRoute, Link } from "@tanstack/react-router";
import { Button } from "@/components/ui/button";
import { Wallet, UserPlus, ArrowRight, Banknote, ListChecks } from "lucide-react";
import { useState } from "react";

export const Route = createFileRoute("/setup")({ component: Setup });

const steps = [
  {
    title: "Create your account",
    description: "Set up your admin credentials to secure your data.",
    icon: UserPlus,
    action: { label: "Register", to: "/login" },
  },
  {
    title: "Add your first account",
    description: "Link a wallet, bank account, or cash balance to start tracking.",
    icon: Banknote,
    action: { label: "Sign in to continue", to: "/login" },
  },
  {
    title: "Start tracking",
    description: "Categorize transactions, set budgets, and monitor your finances.",
    icon: ListChecks,
    action: { label: "Go to Dashboard", to: "/login" },
  },
];

function Setup() {
  const [activeStep, setActiveStep] = useState(0);

  return (
    <div className="flex min-h-screen items-center justify-center bg-zinc-50 p-4 dark:bg-zinc-950">
      <div className="w-full max-w-lg">
        <div className="mb-8 text-center">
          <div className="mb-4 inline-flex h-12 w-12 items-center justify-center rounded-xl bg-zinc-900 text-white dark:bg-zinc-100 dark:text-zinc-900">
            <Wallet className="h-6 w-6" />
          </div>
          <h1 className="text-2xl font-bold text-zinc-900 dark:text-zinc-50">Set up OwnPocket</h1>
          <p className="mt-1 text-zinc-600 dark:text-zinc-400">Follow the steps to get started.</p>
        </div>

        <div className="space-y-4">
          {steps.map((step, i) => {
            const isActive = activeStep === i;
            const isDone = i < activeStep;

            return (
              <button
                key={i}
                type="button"
                onClick={() => setActiveStep(i)}
                className={`w-full rounded-xl border p-5 text-left transition-all ${
                  isActive
                    ? "border-zinc-300 bg-white shadow-sm dark:border-zinc-700 dark:bg-zinc-900"
                    : "border-transparent bg-zinc-100/50 hover:bg-zinc-100 dark:bg-zinc-800/30 dark:hover:bg-zinc-800/50"
                }`}
              >
                <div className="flex items-start gap-4">
                  <span
                    className={`flex h-8 w-8 shrink-0 items-center justify-center rounded-lg text-sm font-semibold ${
                      isDone
                        ? "bg-emerald-500 text-white"
                        : isActive
                          ? "bg-zinc-900 text-white dark:bg-zinc-100 dark:text-zinc-900"
                          : "bg-zinc-200 text-zinc-500 dark:bg-zinc-700 dark:text-zinc-400"
                    }`}
                  >
                    {isDone ? "✓" : i + 1}
                  </span>
                  <div className="flex-1">
                    <div className="flex items-center justify-between">
                      <h3
                        className={`font-semibold ${
                          isActive
                            ? "text-zinc-900 dark:text-zinc-50"
                            : "text-zinc-600 dark:text-zinc-400"
                        }`}
                      >
                        {step.title}
                      </h3>
                      <step.icon
                        className={`h-4 w-4 ${
                          isActive
                            ? "text-zinc-500 dark:text-zinc-400"
                            : "text-zinc-300 dark:text-zinc-600"
                        }`}
                      />
                    </div>
                    <p className="mt-0.5 text-sm text-zinc-500 dark:text-zinc-400">
                      {step.description}
                    </p>
                    {isActive && (
                      <div className="mt-4">
                        <Link to={step.action.to}>
                          <Button size="sm">
                            {step.action.label}
                            <ArrowRight className="ml-1.5 h-3.5 w-3.5" />
                          </Button>
                        </Link>
                      </div>
                    )}
                  </div>
                </div>
              </button>
            );
          })}
        </div>

        <p className="mt-8 text-center text-sm text-zinc-400 dark:text-zinc-500">
          Already have an account?{" "}
          <Link
            to="/login"
            className="font-medium text-zinc-900 underline underline-offset-2 hover:text-zinc-700 dark:text-zinc-300 dark:hover:text-zinc-100"
          >
            Sign in
          </Link>
        </p>
      </div>
    </div>
  );
}

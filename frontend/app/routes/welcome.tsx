import { useEffect, useState } from "react";
import { useNavigate, Link } from "react-router";
import { Button } from "@/components/ui/button";
import { useAuthStore } from "@/stores/auth-store";
import { useThemeStore } from "@/stores/theme-store";
import {
  Wallet,
  ArrowRight,
  Database,
  CheckCircle2,
  XCircle,
  Loader2,
  Sun,
  Moon,
} from "lucide-react";

const API_BASE = import.meta.env.VITE_API_URL ?? "/api/v1";

export function meta() {
  return [{ title: "OwnPocket - Welcome" }];
}

export default function Welcome() {
  const navigate = useNavigate();
  const [status, setStatus] = useState<"loading" | "online" | "offline">("loading");
  const theme = useThemeStore((s) => s.theme);
  const toggleTheme = useThemeStore((s) => s.toggleTheme);

  useEffect(() => {
    const token = useAuthStore.getState().token;
    if (token) {
      localStorage.setItem("ownpocket-has-visited", "true");
      navigate("/", { replace: true });
    }
  }, [navigate]);

  useEffect(() => {
    const controller = new AbortController();
    fetch(`${API_BASE}/health`, { signal: controller.signal })
      .then((res) => setStatus(res.ok ? "online" : "offline"))
      .catch(() => setStatus("offline"));
    return () => controller.abort();
  }, []);

  return (
    <div className="flex min-h-screen items-center justify-center bg-zinc-50 p-4 font-sans dark:bg-zinc-950">
      <div className="absolute inset-0 bg-[linear-gradient(to_right,#e5e7eb_1px,transparent_1px),linear-gradient(to_bottom,#e5e7eb_1px,transparent_1px)] bg-[size:4rem_4rem] [mask-image:radial-gradient(ellipse_60%_50%_at_50%_50%,#000_70%,transparent_100%)] opacity-20 dark:bg-[linear-gradient(to_right,#27272a_1px,transparent_1px),linear-gradient(to_bottom,#27272a_1px,transparent_1px)]" />

      <button
        type="button"
        onClick={toggleTheme}
        className="fixed right-5 top-5 z-50 flex h-9 w-9 items-center justify-center rounded-xl border border-zinc-200 bg-white text-zinc-500 shadow-sm transition-colors hover:bg-zinc-100 hover:text-zinc-700 dark:border-zinc-700 dark:bg-zinc-800 dark:text-zinc-400 dark:hover:bg-zinc-700 dark:hover:text-zinc-200"
      >
        {theme === "dark" ? <Sun className="h-4 w-4" /> : <Moon className="h-4 w-4" />}
      </button>

      <div className="relative w-full max-w-lg rounded-2xl border border-zinc-200 bg-white/80 p-8 shadow-sm backdrop-blur-xl dark:border-zinc-800 dark:bg-zinc-900/80">
        <div className="flex flex-col items-center text-center">
          <div className="mb-6 flex h-16 w-16 items-center justify-center rounded-2xl bg-zinc-900 text-white shadow-inner dark:bg-zinc-100 dark:text-zinc-900">
            <Wallet className="h-8 w-8" />
          </div>

          <h1 className="text-3xl font-bold tracking-tight text-zinc-900 dark:text-zinc-50">
            Welcome to OwnPocket
          </h1>
          <p className="mt-2 text-zinc-600 dark:text-zinc-400">
            Your self-hosted personal finance manager is up and running.
          </p>
        </div>

        <div className="mt-8 rounded-xl border border-zinc-100 bg-zinc-50/50 p-4 dark:border-zinc-800/50 dark:bg-zinc-950/50">
          <ul className="space-y-3 text-sm">
            <li className="flex items-center gap-3 text-zinc-700 dark:text-zinc-300">
              {status === "loading" ? (
                <Loader2 className="h-4 w-4 animate-spin text-zinc-400" />
              ) : status === "online" ? (
                <CheckCircle2 className="h-4 w-4 text-emerald-500" />
              ) : (
                <XCircle className="h-4 w-4 text-red-500" />
              )}
              System
              <span className="text-zinc-400 dark:text-zinc-500">
                {status === "loading" ? "Checking..." : status === "online" ? "Online" : "Offline"}
              </span>
            </li>
            <li className="flex items-center gap-3 text-zinc-700 dark:text-zinc-300">
              <Database className="h-4 w-4 text-zinc-400" />
              SQLite Database
              <span className="text-zinc-400 dark:text-zinc-500">
                {status === "loading"
                  ? "Connecting..."
                  : status === "online"
                    ? "Connected"
                    : "Disconnected"}
              </span>
            </li>
          </ul>
        </div>

        <div className="mt-8 flex flex-col gap-3">
          <Link to="/login" className="w-full">
            <Button size="lg" className="w-full justify-between rounded-xl text-base">
              Access your pocket
              <ArrowRight className="h-4 w-4 opacity-70" />
            </Button>
          </Link>
        </div>
      </div>
    </div>
  );
}

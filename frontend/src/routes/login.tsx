import { createFileRoute, useRouter, Link } from "@tanstack/react-router";
import { useLogin } from "@/hooks/use-auth";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Wallet, ArrowRight, Sun, Moon } from "lucide-react";
import { useState, type FormEvent } from "react";
import { useThemeStore } from "@/stores/theme-store";

export const Route = createFileRoute("/login")({ component: LoginPage });

function LoginPage() {
  const router = useRouter();
  const login = useLogin();
  const theme = useThemeStore((s) => s.theme);
  const toggleTheme = useThemeStore((s) => s.toggleTheme);
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();
    try {
      await login.mutateAsync({ username, password });
      router.navigate({ to: "/" });
    } catch {
      // error displayed in form
    }
  };

  return (
    <div className="flex min-h-screen items-center justify-center bg-zinc-50 p-4 dark:bg-zinc-950">
      <button
        type="button"
        onClick={toggleTheme}
        className="fixed right-5 top-5 z-50 flex h-9 w-9 items-center justify-center rounded-xl border border-zinc-200 bg-white text-zinc-500 shadow-sm transition-colors hover:bg-zinc-100 hover:text-zinc-700 dark:border-zinc-700 dark:bg-zinc-800 dark:text-zinc-400 dark:hover:bg-zinc-700 dark:hover:text-zinc-200"
      >
        {theme === "dark" ? <Sun className="h-4 w-4" /> : <Moon className="h-4 w-4" />}
      </button>

      <div className="relative w-full max-w-sm rounded-2xl border border-zinc-200 bg-white/80 p-8 shadow-sm backdrop-blur-xl dark:border-zinc-800 dark:bg-zinc-900/80">
        <div className="mb-6 flex flex-col items-center text-center">
          <div className="mb-4 flex h-12 w-12 items-center justify-center rounded-xl bg-zinc-900 text-white shadow-inner dark:bg-zinc-100 dark:text-zinc-900">
            <Wallet className="h-6 w-6" />
          </div>
          <h1 className="text-xl font-bold text-zinc-900 dark:text-zinc-50">
            Sign in to OwnPocket
          </h1>
        </div>

        <form onSubmit={handleSubmit} className="flex flex-col gap-4">
          <Input
            label="Username"
            id="username"
            type="text"
            value={username}
            onChange={(e) => setUsername(e.target.value)}
            required
          />
          <Input
            label="Password"
            id="password"
            type="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            required
          />
          {login.error && (
            <p className="text-sm text-red-600 dark:text-red-400">
              {login.error instanceof Error ? login.error.message : "Login failed"}
            </p>
          )}
          <Button type="submit" className="w-full justify-between" disabled={login.isPending}>
            {login.isPending ? "Signing in..." : "Sign in"}
            <ArrowRight className="h-4 w-4 opacity-70" />
          </Button>
        </form>

        <p className="mt-6 text-center text-sm text-zinc-500 dark:text-zinc-400">
          Don&apos;t have an account?{" "}
          <Link
            to="/setup"
            className="font-medium text-zinc-900 underline underline-offset-2 hover:text-zinc-700 dark:text-zinc-300 dark:hover:text-zinc-100"
          >
            Set up OwnPocket
          </Link>
        </p>
      </div>
    </div>
  );
}

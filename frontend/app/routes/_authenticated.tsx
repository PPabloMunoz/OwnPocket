import { useEffect, useState } from "react";
import { Outlet, useNavigate } from "react-router";
import { AppShell } from "@/components/layout";
import { useAuthStore } from "@/stores/auth-store";

export default function AuthenticatedLayout() {
  const navigate = useNavigate();
  const [checking, setChecking] = useState(true);

  useEffect(() => {
    const hasVisited = localStorage.getItem("ownpocket-has-visited");
    if (!hasVisited) {
      localStorage.setItem("ownpocket-has-visited", "true");
      navigate("/welcome", { replace: true });
      return;
    }
    const token = useAuthStore.getState().token;
    if (!token) {
      navigate("/login", { replace: true });
      return;
    }
    setChecking(false);
  }, [navigate]);

  if (checking) {
    return (
      <div className="flex min-h-screen items-center justify-center bg-zinc-100 dark:bg-zinc-900" />
    );
  }

  return (
    <AppShell>
      <Outlet />
    </AppShell>
  );
}

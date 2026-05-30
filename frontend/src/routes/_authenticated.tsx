import { createFileRoute, Outlet, redirect } from "@tanstack/react-router";
import { AppShell } from "@/components/layout";
import { useAuthStore } from "@/stores/auth-store";

export const Route = createFileRoute("/_authenticated")({
  beforeLoad: () => {
    if (typeof window !== "undefined") {
      const hasVisited = localStorage.getItem("ownpocket-has-visited");
      const token = useAuthStore.getState().token;
      console.log("visited: ", hasVisited);
      console.log("token: ", token);
      if (!hasVisited) {
        localStorage.setItem("ownpocket-has-visited", "true");
        throw redirect({ to: "/welcome" });
      }
      if (!token) {
        throw redirect({ to: "/login" });
      }
    }
  },
  component: AuthenticatedLayout,
});

function AuthenticatedLayout() {
  return (
    <AppShell>
      <Outlet />
    </AppShell>
  );
}

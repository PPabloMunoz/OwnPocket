import { lazy, Suspense, useEffect, useState } from "react";
import {
  Links,
  Meta,
  Outlet,
  Scripts,
  ScrollRestoration,
  isRouteErrorResponse,
  useRouteError,
} from "react-router";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { Frown, AlertTriangle, Home, ArrowLeft } from "lucide-react";

const ReactQueryDevtoolsProduction = lazy(() =>
  import("@tanstack/react-query-devtools/build/modern/production.js").then((d) => ({
    default: d.ReactQueryDevtools,
  })),
);
import { useThemeStore } from "@/stores/theme-store";
import appCss from "./app.css?url";

const queryClient = new QueryClient();

export function meta() {
  return [{ title: "OwnPocket" }, { name: "description", content: "Personal finance manager" }];
}

export function links() {
  return [{ rel: "stylesheet", href: appCss }];
}

export function Layout({ children }: { children: React.ReactNode }) {
  const theme = useThemeStore((s) => s.theme);
  const [transitioning, setTransitioning] = useState(false);

  useEffect(() => {
    requestAnimationFrame(() => setTransitioning(true));
  }, []);

  return (
    <html lang="en" data-theme={theme} className={transitioning ? "theme-transition" : ""}>
      <head>
        <meta charSet="utf-8" />
        <meta name="viewport" content="width=device-width, initial-scale=1" />
        <Meta />
        <Links />
      </head>
      <body>
        {children}
        <ScrollRestoration />
        <Scripts />
      </body>
    </html>
  );
}

export default function App() {
  const [showDevtools, setShowDevtools] = useState(false);

  useEffect(() => {
    // @ts-ignore
    window.toggleDevtools = () => setShowDevtools((prev) => !prev);
  }, []);

  return (
    <QueryClientProvider client={queryClient}>
      <Outlet />
      {showDevtools && (
        <Suspense fallback={null}>
          <ReactQueryDevtoolsProduction />
        </Suspense>
      )}
    </QueryClientProvider>
  );
}

export function ErrorBoundary() {
  const error = useRouteError();
  let details = "An unexpected error occurred.";
  let stack: string | undefined;

  if (isRouteErrorResponse(error)) {
    details =
      error.status === 404 ? "The requested page could not be found." : error.statusText || details;
  } else if (import.meta.env.DEV && error && error instanceof Error) {
    details = error.message;
    stack = error.stack;
  }

  const is404 = isRouteErrorResponse(error) && error.status === 404;

  return (
    <div className="flex min-h-screen items-center justify-center bg-zinc-50 p-4 dark:bg-zinc-950">
      <div className="relative w-full max-w-md text-center">
        <div className="mb-6 flex justify-center">
          <div
            className={`flex h-20 w-20 items-center justify-center rounded-3xl ${
              is404 ? "bg-zinc-100 dark:bg-zinc-800" : "bg-red-50 dark:bg-red-950/50"
            }`}
          >
            {is404 ? (
              <Frown className="h-10 w-10 text-zinc-400 dark:text-zinc-500" />
            ) : (
              <AlertTriangle className="h-10 w-10 text-red-400 dark:text-red-500" />
            )}
          </div>
        </div>
        {is404 ? (
          <>
            <h1 className="text-6xl font-bold tracking-tight text-zinc-200 dark:text-zinc-800">
              404
            </h1>
            <h2 className="-mt-2 text-xl font-semibold text-zinc-900 dark:text-zinc-50">
              Page not found
            </h2>
            <p className="mt-2 text-sm text-zinc-500 dark:text-zinc-400">
              The page you&apos;re looking for doesn&apos;t exist or has been moved.
            </p>
          </>
        ) : (
          <>
            <h1 className="text-2xl font-bold tracking-tight text-zinc-900 dark:text-zinc-50">
              Something went wrong
            </h1>
            <p className="mt-2 text-sm text-zinc-500 dark:text-zinc-400">
              An unexpected error occurred. Please try again.
            </p>
            {details && (
              <p className="mt-4 rounded-xl border border-red-100 bg-red-50/50 px-4 py-3 text-left text-xs font-mono text-red-700 dark:border-red-900/50 dark:bg-red-950/30 dark:text-red-400">
                {details}
              </p>
            )}
            {stack && (
              <pre className="mt-2 w-full overflow-x-auto rounded-xl border border-red-100 bg-red-50/50 p-4 text-left text-xs font-mono text-red-700 dark:border-red-900/50 dark:bg-red-950/30 dark:text-red-400">
                <code>{stack}</code>
              </pre>
            )}
          </>
        )}
        <div className="mt-8 flex justify-center gap-3">
          <a
            href="/"
            className="inline-flex items-center gap-2 rounded-xl bg-zinc-900 px-5 py-2.5 text-sm font-medium text-white transition-colors hover:bg-zinc-800 dark:bg-zinc-100 dark:text-zinc-900 dark:hover:bg-zinc-200"
          >
            <Home className="h-4 w-4" />
            Go home
          </a>
          <button
            type="button"
            onClick={() => window.history.back()}
            className="inline-flex items-center gap-2 rounded-xl border border-zinc-200 bg-white px-5 py-2.5 text-sm font-medium text-zinc-700 transition-colors hover:bg-zinc-50 dark:border-zinc-700 dark:bg-zinc-800 dark:text-zinc-300 dark:hover:bg-zinc-700"
          >
            <ArrowLeft className="h-4 w-4" />
            Go back
          </button>
        </div>
      </div>
    </div>
  );
}

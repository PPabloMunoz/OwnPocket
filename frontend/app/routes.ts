import {
  type RouteConfig,
  route,
  index,
  layout,
} from "@react-router/dev/routes";

export default [
  route("welcome", "routes/welcome.tsx"),
  route("login", "routes/login.tsx"),
  route("register", "routes/register.tsx"),
  route("setup", "routes/setup.tsx"),
  layout("routes/_authenticated.tsx", [
    index("routes/dashboard.tsx"),
    route("accounts", "routes/accounts.tsx"),
    route("transactions", "routes/transactions.tsx"),
    route("budgets", "routes/budgets.tsx"),
    route("categories", "routes/categories.tsx"),
  ]),
] satisfies RouteConfig;

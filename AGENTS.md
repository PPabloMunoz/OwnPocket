# AGENTS.md - AI Assistant Development Guide

This document is written specifically for AI coding agents (Grok, Claude, Cursor, Windsurf, etc.) to help you understand, maintain, and extend **OwnPocket** — a lightweight, self-hosted personal finance manager.

---

## Project Overview

**OwnPocket** is a minimalist, performant, self-hosted income/expenses/budget tracker designed for easy deployment on Proxmox LXC via [community-scripts.org](https://community-scripts.org).

**Core Goals:**

- Single binary deployment (Go + embedded frontend)
- Very low resource usage (ideal for LXC)
- Excellent developer experience for both humans and AI agents
- SQLite-first with all monetary values in **cents** (integers)
- Clean architecture

---

## Technology Stack

| Layer | Technology | Reason |
| :--- | :--- | :--- |
| Backend | Go 1.23+ + Gin | Performance, single binary |
| Database | SQLite + GORM | Zero-config, embedded |
| Frontend | Vite + React 19 + TypeScript + Tailwind | Fast development, modern UI |
| Styling | Tailwind CSS v4 + DaisyUI v5 + lucide-react icons | Utility-first, accessible |
| State | TanStack Query (server) + Zustand (client) | Server + client state |
| Routing | React Router v7 (file-based via `routes.ts`) | Type-safe, nested layouts, SPA mode |
| Linting | oxlint + oxfmt (Rust-based) | Lightning-fast linting & formatting |
| Auth | JWT + bcrypt | Simple & secure |
| Money | Custom `model.Amount` (int64 cents) | Precision safety |

---

## Project Structure

```bash
OwnPocket/
├── backend/
│   ├── cmd/server/main.go              # Application entrypoint
│   ├── internal/
│   │   ├── config/                     # Configuration & env
│   │   │   ├── config.go
│   │   │   └── database.go
│   │   ├── handler/                    # HTTP handlers (Gin)
│   │   │   ├── handler.go              # Route registration
│   │   │   ├── frontend.go             # Embedded frontend handler
│   │   │   ├── account_handler.go
│   │   │   ├── auth_handler.go
│   │   │   ├── budget_handler.go
│   │   │   ├── category_handler.go
│   │   │   ├── dashboard_handler.go
│   │   │   └── transaction_handler.go
│   │   ├── service/                    # Business logic
│   │   │   ├── service.go              # Service base/deps
│   │   │   ├── account_service.go
│   │   │   ├── auth_service.go
│   │   │   ├── budget_service.go
│   │   │   ├── category_service.go
│   │   │   ├── dashboard_service.go
│   │   │   └── transaction_service.go
│   │   ├── model/                      # GORM models + types
│   │   │   ├── models.go
│   │   │   └── amount.go              # Custom Amount (int64 cents)
│   │   ├── middleware/                 # Auth, CORS, etc.
│   │   │   └── middleware.go
│   │   └── utils/                      # Shared helpers
│   │       └── response.go
│   │
│   ├── .air.toml                       # Air hot-reload config
│   ├── .env / .env.example
│   └── go.mod
│
├── frontend/
│   ├── app/
│   │   ├── app.css                    # Tailwind v4 + DaisyUI entry
│   │   ├── root.tsx                   # Root layout (QueryClientProvider, meta, links)
│   │   ├── routes.ts                  # React Router v7 route config (flat file-based)
│   │   ├── components/
│   │   │   ├── ui/                    # Base UI primitives (button, input, card, badge, dialog)
│   │   │   └── layout/               # app-shell, navbar
│   │   ├── hooks/                     # Custom hooks (use-auth)
│   │   ├── lib/                       # API client, query keys, utils
│   │   ├── stores/                    # Zustand stores (auth-store, theme-store)
│   │   ├── types/                     # TypeScript types (one file per domain)
│   │   └── routes/                    # Route page components (flat)
│   │       ├── _authenticated.tsx     # Protected layout
│   │       ├── login.tsx
│   │       ├── register.tsx
│   │       ├── setup.tsx              # Setup wizard
│   │       ├── welcome.tsx            # Welcome splash
│   │       ├── dashboard.tsx
│   │       ├── accounts.tsx
│   │       ├── transactions.tsx
│   │       ├── budgets.tsx
│   │       └── categories.tsx
│   ├── .react-router/                 # Auto-generated types by React Router
│   ├── package.json
│   ├── vite.config.ts
│   ├── react-router.config.ts
│   └── tsconfig.json
│
├── data/                               # Runtime data (gitignore)
│   ├── app.db
│   ├── test_data.sql                  # Sample data for development
│   └── migrations/
│
├── proxmox/                            # Proxmox LXC installation scripts
│   ├── ct/                             # Container scripts
│   └── install/                        # Installation scripts
│
├── bin/                                # Compiled binaries
│   └── app                            # Single executable
├── justfile
├── Dockerfile
├── docker-compose.yml
├── LICENSE.md
└── README.md
```

---

## Architecture Principles (Follow These)

1. **Clean Architecture** — Handlers → Service → Repository (via GORM)
2. **All money = `model.Amount` (int64 cents)**
3. **User-scoped queries** — Always filter by `user_id`
4. **Database transactions** for balance updates
5. **Feature-first development** — Add small, complete features
6. **Single responsibility** — One file per major entity when possible

---

## Coding Standards

### Go Standards

- Use **internal/** for all non-exported code
- Prefer explicit error handling
- Always validate input in handlers
- Use `ShouldBindJSON` with struct tags
- Return consistent JSON format: `{ "data": ..., "error": null }` / `{ "data": null, "error": "..." }` using `utils.Success()` and `utils.Error()`
- Comment complex business logic
- Keep services focused (one file per domain when possible)

### TypeScript / React Standards

- Use TypeScript strictly
- Prefer functional components + hooks
- Use `zustand` for client state (persisted auth, UI preferences)
- Use TanStack Query for server state (API data fetching/caching)
- Component filenames: `kebab-case.tsx`
- API calls centralized in `app/lib/api.ts`
- Types mirroring backend models go in `app/types/` (one file per domain)
- Use `@/` path alias for `app/` imports (e.g. `@/lib/api`, `@/components/ui/button`)
- Routes follow React Router v7 file-based convention in `app/routes.ts`
- UI primitives go in `app/components/ui/`, layout components in `app/components/layout/`
- Prefer `import type` for type-only imports (required by `verbatimModuleSyntax`)
- Custom hooks go in `app/hooks/`

---

## Key Types

### `model.Amount`

```go
type Amount int64

func NewAmountFromCents(cents int64) Amount
func NewAmountFromFloat(dollars float64) Amount
func (a Amount) Cents() int64
func (a Amount) ToFloat() float64
```

**Never** use `float64` for money in the backend.

---

## How to Add a New Feature (Step-by-step)

### Example: Adding "Tags Management"

1. **Model** (`internal/model/models.go`)
   - Add/update structs if needed

2. **Service** (`internal/service/tag_service.go`)
   - Create CRUD methods
   - Write tests in `tag_service_test.go`

3. **Handler** (`internal/handler/tag_handler.go`)
   - Add routes and handlers
   - Write tests in `tag_handler_test.go`

4. **Routes** (`internal/handler/handler.go`)
   - Register new routes under protected group

5. **Frontend** (if applicable)
   - Add types in `app/types/`
   - Create components in `app/components/`
   - Add API calls in `app/lib/api.ts`
   - Update route in `app/routes.ts` if needed
   - Create page component in `app/routes/`

6. **Test** the full flow

---

## Important Business Rules

- **Transactions**:
  - `amount` is **always positive**
  - Direction is determined by `type` (`income`/`expense`/`transfer`)
  - Balance updates must be atomic (use `db.Transaction`)

- **Transfers**:
  - Must update both source and destination accounts

- **Budgets**:
  - Monthly envelopes (`period` = `YYYY-MM`)

- **Soft Delete**:
  - Use `is_active` on accounts instead of hard delete

---

## Common Tasks for AI Agents

### When asked to implement something:

1. Read relevant existing files first
2. Follow the same pattern as similar features
3. Update `AGENTS.md` if you introduce new conventions
4. Keep changes minimal and focused
5. Add proper error handling and validation
6. Make sure it works with the existing auth system

### When reviewing code:

- Check for float usage with money
- Verify user-scoped queries
- Ensure database transactions for balance changes
- Validate input
- Check consistency with existing style

---

## Development Workflow

```bash
# Backend only
just dev-backend

# Frontend only
just dev-frontend

# Full stack
just dev

# Build single binary
just build-local
```

---

## Deployment Notes

- Single static binary (Go embeds frontend)
- SQLite file mounted as volume
- Nginx/Caddy recommended in front for production
- Designed for Proxmox LXC (low memory ~150-300MB)

---

## Future Extension Points

- Multi-currency support
- Split transactions
- Receipt attachments
- Rules engine
- Import from CSV / OFX / QIF
- Multi-user with proper permissions
- Export reports (PDF)

---

## Contribution Guidelines for Agents

- Be explicit about what you're changing
- Keep PRs focused (one feature/fix per change)
- Update this `AGENTS.md` when introducing new patterns
- Prefer composition over inheritance
- Document any non-obvious business logic
- **Proxmox LXC Script**: If tasked with creating or maintaining the installation script for community-scripts.org, follow the [Community Scripts Agent Guide](https://community-scripts.org/docs/contribution/agents).

---

**Last Updated:** May 31, 2026 (added Community Scripts agent guide)

You are now a **OwnPocket expert agent**. Follow this document strictly when helping with this codebase.

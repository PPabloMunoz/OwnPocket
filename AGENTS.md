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

| Layer    | Technology                           | Reason                      |
| -------- | ------------------------------------ | --------------------------- |
| Backend  | Go 1.23+ + Gin                       | Performance, single binary  |
| Database | SQLite + GORM                        | Zero-config, embedded       |
| Frontend | Vite + React + TypeScript + Tailwind | Fast development, modern UI |
| Styling  | Tailwind + shadcn/ui                 | Beautiful + accessible      |
| Auth     | JWT + bcrypt                         | Simple & secure             |
| Money    | Custom `model.Amount` (int64 cents)  | Precision safety            |

---

## Project Structure

```bash
OwnPocket/
├── backend/
│   ├── cmd/server/main.go              # Application entrypoint
│   ├── internal/
│   │   ├── config/                     # Configuration & env
│   │   ├── handler/                    # HTTP handlers (Gin)
│   │   ├── service/                    # Business logic
│   │   ├── model/                      # GORM models + types
│   │   ├── middleware/                 # Auth, CORS, etc.
│   │   └── utils/                      # Shared helpers (responses, etc.)
│   ├── migrations/                     # SQL migration files
│   └── go.mod
│
├── data/                               # Runtime data (gitignore)
│   └── app.db
│
├── justfile
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
- Use `zustand` for state
- Use TanStack Query for server state
- Component filenames: `PascalCase.tsx`
- API calls centralized in `src/lib/api.ts`

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

3. **Handler** (`internal/handler/tag_handler.go`)
   - Add routes and handlers

4. **Routes** (`internal/handler/handler.go`)
   - Register new routes under protected group

5. **Frontend**
   - Add types in `src/types/`
   - Create components
   - Add API calls in `api.ts`
   - Update relevant pages

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

---

**Last Updated:** May 29, 2026 (updated entry point, build system, project structure, JSON response format)

You are now a **OwnPocket expert agent**. Follow this document strictly when helping with this codebase.


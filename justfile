# lite-finance tasks runner
# Uses 'just' to automate backend, frontend, and Docker workflows

set shell := ["sh", "-c"]

# Display available commands
default:
    @just --list

# ---
# Development Commands
# ---

# Run the backend in development mode (with hot-reloading if air is installed, otherwise standard go run)
dev-backend:
    @echo "🚀 Starting Go backend..."
    cd backend && if command -v air >/dev/null 2>&1; then air; else go run cmd/server/main.go; fi

# Run the frontend development server
dev-frontend:
    @echo "⚡ Starting Vite frontend..."
    cd frontend && pnpm run dev

# Run both backend and frontend concurrently (requires 'concurrently' npm package or similar, or runs in background)
dev:
    @echo "🛠️ Starting full-stack development environment..."
    just -j 2 dev-backend dev-frontend

# ---
# Installation & Setup
# ---

# Install all dependencies for both backend and frontend
setup:
    @echo "📦 Installing frontend dependencies..."
    cd frontend && pnpm install
    @echo "📥 Tidy Go modules..."
    cd backend && go mod tidy && go mod vendor
    @echo "📁 Creating data directory for SQLite..."
    mkdir -p data
    @echo "✅ Setup complete."

# ---
# Testing & Linting
# ---

# Run tests for backend and frontend
test:
    @echo "🧪 Running backend tests..."
    cd backend && go test ./...
    @echo "🧪 Running frontend tests..."
    cd frontend && pnpm run test -- --watchAll=false

# Lint the codebase
lint:
    @echo "🧹 Linting backend..."
    cd backend && if command -v golangci-lint >/dev/null 2>&1; then golangci-lint run; else go fmt ./...; fi
    @echo "🧹 Linting frontend..."
    cd frontend && npm run lint

# ---
# Build Commands
# ---

# Build the frontend production assets
build-frontend:
    @echo "📦 Building frontend production bundles..."
    cd frontend && npm run build

# Build the backend binary
build-backend:
    @echo "🏗️ Building Go backend binary..."
    cd backend && go build -o ../bin/app cmd/server/main.go

# Build both frontend and backend locally
build-local: build-frontend build-backend

# ---
# Docker & Deployment
# ---

# Spin up the entire stack using Docker Compose
up:
    @echo "🐋 Starting production containers..."
    docker-compose up -d

# Stop the Docker Compose stack
down:
    @echo "🛑 Stopping containers..."
    docker-compose down

# Rebuild and restart the Docker environment
restart:
    @echo "🔄 Rebuilding and restarting containers..."
    docker-compose down
    docker-compose up -d --build

# View real-time logs from Docker containers
logs:
    docker-compose logs -f

# ---
# Database Utilities
# ---

# Quick access to look into the local SQLite database (requires sqlite3 CLI tool)
db-shell:
    @if [ -f "data/app.db" ]; then sqlite3 data/app.db; else echo "❌ Database file 'data/app.db' not found. Run the app first."; fi

# Clean build artifacts, node_modules, and binaries (keeps the DB intact)
clean:
    @echo "🧹 Cleaning up build artifacts..."
    rm -rf bin/
    rm -rf frontend/dist
    @echo "✨ Clean complete."

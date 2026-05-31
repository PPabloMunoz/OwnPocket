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

# Lint the codebase
lint:
    @echo "🧹 Linting backend..."
    cd backend && if command -v golangci-lint >/dev/null 2>&1; then golangci-lint run; else go fmt ./...; fi
    @echo "🧹 Linting frontend..."
    cd frontend && pnpm lint

# ---
# Build Commands
# ---

# Build the frontend production assets
build-frontend:
    @echo "📦 Building frontend production bundles..."
    cd frontend && pnpm run build
    @echo "🚚 Copying frontend assets to backend..."
    rm -rf backend/internal/handler/static/*
    cp -rv frontend/build/client/* backend/internal/handler/static/

# Build the backend binary
build-backend:
    @echo "🏗️ Building Go backend binary..."
    cd backend && go build -o ../bin/app cmd/server/main.go

# Build both frontend and backend locally
build-local: build-frontend build-backend

# Prepare a release (builds everything and shows the binary path)
release: build-local
    @echo "✅ Release ready at ./bin/app"
    @ls -lh ./bin/app

# ---
# Docker & Deployment
# ---

# Build the Docker image
docker-build:
    @echo "🐳 Building Docker image..."
    docker build -t ownpocket .

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

# Create a new git tag and push it (triggers the release workflow)
tag version:
    @echo "🏷️ Creating tag {{version}}..."
    git tag -a {{version}} -m "Release {{version}}"
    git push origin {{version}}
    @echo "🚀 Tag pushed! The release workflow should start shortly on GitHub."

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
    rm -rf frontend/build
    rm -rf backend/internal/handler/static/*
    @echo "✨ Clean complete."

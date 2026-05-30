# Stage 1: Build Frontend
FROM node:22-alpine AS frontend-builder
WORKDIR /app/frontend
COPY frontend/package.json frontend/pnpm-lock.yaml ./
RUN npm install -g pnpm && pnpm install
COPY frontend/ ./
RUN pnpm run build

# Stage 2: Build Backend
FROM golang:alpine AS backend-builder
WORKDIR /app/backend
# Install build essentials if needed
RUN apk add --no-cache gcc musl-dev
COPY backend/go.mod backend/go.sum ./
RUN go mod download
COPY backend/ ./
# Copy built frontend assets from Stage 1
COPY --from=frontend-builder /app/frontend/build/client ./internal/handler/static
# Build the binary
RUN CGO_ENABLED=1 GOOS=linux go build -o /app/bin/app cmd/server/main.go

# Stage 3: Final Image
FROM alpine:latest
WORKDIR /app
# Install dependencies for CGO (sqlite)
RUN apk add --no-cache ca-certificates libc6-compat
# Copy the binary from the backend-builder
COPY --from=backend-builder /app/bin/app .
# Create data directory
RUN mkdir -p data
# Set environment variables
ENV PORT=8080
ENV DB_PATH="/app/data/app.db"
ENV GIN_MODE=release
# Expose the port
EXPOSE 8080
# Run the application
CMD ["./app"]

# ==========================================
# STAGE 1: Frontend Dependency Caching
# ==========================================
FROM node:26-alpine3.22 AS frontend-deps
WORKDIR /app/frontend
RUN apk add --no-cache libc6-compat
RUN npm install -g pnpm
COPY frontend/package.json frontend/pnpm-lock.yaml frontend/pnpm-workspace.yaml ./
RUN pnpm approve-builds esbuild && pnpm install --frozen-lockfile

# ==========================================
# STAGE 2: Frontend Build
# ==========================================
FROM frontend-deps AS frontend-builder
WORKDIR /app/frontend
COPY frontend/ ./
RUN pnpm run build

# ==========================================
# STAGE 3: Backend Dependency Caching
# ==========================================
FROM golang:1.26-alpine AS backend-deps
WORKDIR /app/backend
RUN apk add --no-cache gcc musl-dev
COPY backend/go.mod backend/go.sum ./
RUN go mod download

# ==========================================
# STAGE 4: Backend Build
# ==========================================
FROM backend-deps AS backend-builder
WORKDIR /app/backend
# Copy the source code
COPY backend/ ./
# Copy built frontend assets *ONLY* at the last possible second
COPY --from=frontend-builder /app/frontend/build/client ./internal/handler/static

# Leverage Go's compiler cache for CGO builds
RUN --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=1 GOOS=linux go build -ldflags="-s -w" -o /app/bin/app cmd/server/main.go

# ==========================================
# STAGE 5: Final Production Image
# ==========================================
FROM alpine:3.22
WORKDIR /app

# Combine RUN commands to keep layers minimal
RUN apk add --no-cache ca-certificates libc6-compat && \
    mkdir -p data && \
    addgroup -S appgroup && adduser -S appuser -G appgroup && \
    chown -R appuser:appgroup /app

COPY --from=backend-builder --chown=appuser:appgroup /app/bin/app .

USER appuser

ENV PORT=8080
ENV DB_PATH="/app/data/app.db"
ENV GIN_MODE=release

EXPOSE 8080

CMD ["./app"]

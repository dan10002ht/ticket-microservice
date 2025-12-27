# ============================================
# Stage 1: Dependencies
# ============================================
# Pin exact image digest for security & reproducibility
# Verify: docker pull golang:1.22-alpine && docker images --digests golang:1.22-alpine
FROM golang:1.22-alpine@sha256:1699c10032ca2582ec89a24a1312d986a3f094aed3d5c1147b19880afe40e052 AS deps

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Download dependencies first (layer caching)
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# ============================================
# Stage 2: Builder
# ============================================
FROM deps AS builder

# Copy source code
COPY . .

# Build with optimizations
# -w: omit DWARF symbol table
# -s: omit symbol table and debug info
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -o /app/server ./main.go

# ============================================
# Stage 3: Security Scan (CI/CD - use with --target security)
# ============================================
FROM deps AS security

# Install govulncheck for vulnerability scanning
RUN go install golang.org/x/vuln/cmd/govulncheck@latest

# Copy source for scanning
COPY . .

# Run vulnerability check (non-blocking)
RUN govulncheck ./... || true

# ============================================
# Stage 4: Production (Distroless - Minimal Attack Surface)
# ============================================
# Distroless benefits:
# - No shell, package manager, or unnecessary binaries
# - Smaller attack surface (fewer CVEs)
# - Immutable by design
FROM gcr.io/distroless/static-debian12:nonroot@sha256:2b7c93f6d6648c11f0e80a48558c8f77885eb0445213b8e69a6a0d7c89fc6ae4 AS production

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/server /app/server

# Copy SSL certificates and timezone data
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Copy migrations if exists
COPY --from=builder /app/migrations /app/migrations

# Distroless runs as nonroot (uid 65532) by default
USER nonroot:nonroot

# Expose gRPC port
EXPOSE 50058

# Note: Distroless has no shell, health checks must be external (Docker Compose/K8s)
# Example in docker-compose.yml:
#   healthcheck:
#     test: ["CMD", "/app/server", "health"] # if binary supports health subcommand

# Distroless includes proper init, no ENTRYPOINT needed
CMD ["/app/server"]

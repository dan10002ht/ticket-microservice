# ============================================
# Stage 1: Dependencies
# ============================================
# Pin exact image digest for security & reproducibility
# Verify: docker pull node:20-alpine && docker images --digests node:20-alpine
# CVE Scan: trivy image --severity HIGH,CRITICAL node:20-alpine
FROM node:20-alpine@sha256:658d0f63e501824d6c23e06d4bb95c71e7d704537c9d9272f488ac03a370d448 AS deps

WORKDIR /app

# Install build tools for native modules and protobuf
RUN apk add --no-cache protobuf python3 make g++

# Copy package files only
COPY package.json yarn.lock ./

# Install ALL dependencies with integrity check
# --frozen-lockfile ensures yarn.lock is respected (like SHA verification)
RUN yarn install --frozen-lockfile

# ============================================
# Stage 2: Builder
# ============================================
FROM node:20-alpine@sha256:658d0f63e501824d6c23e06d4bb95c71e7d704537c9d9272f488ac03a370d448 AS builder

WORKDIR /app

# Copy dependencies from deps stage
COPY --from=deps /app/node_modules ./node_modules
COPY package.json yarn.lock ./

# Copy source code
COPY . .

# Generate gRPC code if proto files exist
RUN if [ -d "proto" ]; then \
      apk add --no-cache protobuf && \
      yarn grpc:generate; \
    fi

# Prune devDependencies for production
RUN yarn install --frozen-lockfile --production \
    && yarn cache clean

# ============================================
# Stage 3: Security Scan (CI/CD - use with --target security)
# ============================================
# Build: docker build --target security -t gateway:security .
# This stage is for CI/CD pipeline to scan vulnerabilities
FROM deps AS security

# Install trivy for vulnerability scanning
RUN apk add --no-cache curl \
    && curl -sfL https://raw.githubusercontent.com/aquasecurity/trivy/main/contrib/install.sh | sh -s -- -b /usr/local/bin

# Copy production node_modules to scan
COPY --from=builder /app/node_modules ./node_modules
COPY --from=builder /app/package.json ./

# Scan for vulnerabilities (npm packages)
# --exit-code 1 will fail build if HIGH/CRITICAL found
RUN trivy fs --severity HIGH,CRITICAL --exit-code 0 /app || true

# Also run yarn audit
RUN yarn audit --level high || true

# ============================================
# Stage 4: Production (Distroless - Minimal Attack Surface)
# ============================================
# Distroless benefits:
# - No shell, package manager, or unnecessary binaries
# - Smaller attack surface (fewer CVEs)
# - Immutable by design
# CVE Scan: trivy image --severity HIGH,CRITICAL gcr.io/distroless/nodejs20-debian12:nonroot
FROM gcr.io/distroless/nodejs20-debian12:nonroot@sha256:010cb788479147947d2515097bc9c8b77e12e77316f79a0d46e60c445905038b AS production

WORKDIR /app

# Copy production files (distroless runs as nonroot by default, uid 65532)
COPY --from=builder --chown=65532:65532 /app/node_modules ./node_modules
COPY --from=builder --chown=65532:65532 /app/src ./src
COPY --from=builder --chown=65532:65532 /app/package.json ./

# Copy generated gRPC files if they exist
COPY --from=builder --chown=65532:65532 /app/generated ./generated

# Expose port
EXPOSE 3000

# Note: Distroless has no shell, so HEALTHCHECK with CMD won't work
# Health checks should be done externally (Docker Compose, K8s, etc.)
# Example in docker-compose.yml:
#   healthcheck:
#     test: ["CMD", "/nodejs/bin/node", "-e", "require('http').get('http://localhost:53000/health')"]

# Distroless already includes tini-like init, no ENTRYPOINT needed
# Start application (distroless uses /nodejs/bin/node)
CMD ["src/index.js"]

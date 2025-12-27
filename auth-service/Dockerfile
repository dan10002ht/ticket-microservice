# ============================================
# Stage 1: Dependencies
# ============================================
# Pin exact image digest for security & reproducibility
# Verify: docker pull node:20-alpine && docker images --digests node:20-alpine
FROM node:20-alpine@sha256:658d0f63e501824d6c23e06d4bb95c71e7d704537c9d9272f488ac03a370d448 AS deps

WORKDIR /app

# Install build tools for native modules (bcrypt, etc.)
RUN apk add --no-cache python3 make g++

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

# Prune devDependencies for production
RUN yarn install --frozen-lockfile --production \
    && yarn cache clean

# ============================================
# Stage 3: Security Scan (Optional - for CI/CD)
# ============================================
FROM deps AS security

# Run security audit (non-blocking)
RUN yarn audit --level moderate || true

# Note: In CI/CD pipeline, use dedicated tools:
# - Trivy for container scanning
# - Snyk for dependency vulnerabilities
# - npm audit / yarn audit for package audit

# ============================================
# Stage 4: Production
# ============================================
FROM node:20-alpine@sha256:658d0f63e501824d6c23e06d4bb95c71e7d704537c9d9272f488ac03a370d448 AS production

# Add tini for proper signal handling (graceful shutdown)
RUN apk add --no-cache tini

# Create non-root user FIRST (before copying files)
RUN addgroup -g 1001 -S nodejs \
    && adduser -S -u 1001 -G nodejs nodejs

WORKDIR /app

# Create logs directory with correct ownership
RUN mkdir -p logs && chown nodejs:nodejs logs

# Copy production files with correct ownership
COPY --from=builder --chown=nodejs:nodejs /app/node_modules ./node_modules
COPY --from=builder --chown=nodejs:nodejs /app/src ./src
COPY --from=builder --chown=nodejs:nodejs /app/package.json ./

# Copy knexfile and migrations for runtime migrations
COPY --from=builder --chown=nodejs:nodejs /app/knexfile.js ./
COPY --from=builder --chown=nodejs:nodejs /app/migrations ./migrations

# Switch to non-root user
USER nodejs

# Expose gRPC and metrics ports
EXPOSE 50051 9090

# Health check for gRPC service
HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
    CMD node -e "const grpc = require('@grpc/grpc-js'); const client = new grpc.Client('localhost:50051', grpc.credentials.createInsecure()); client.waitForReady(Date.now() + 5000, (err) => { process.exit(err ? 1 : 0); });"

# Use tini as init system for proper signal handling
ENTRYPOINT ["/sbin/tini", "--"]

# Start application
CMD ["node", "src/index.js"]

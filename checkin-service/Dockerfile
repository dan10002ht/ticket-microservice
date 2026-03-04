# ============================================
# Stage 1: Dependencies
# ============================================
FROM golang:1.22-alpine AS deps

WORKDIR /app

RUN apk add --no-cache git ca-certificates tzdata

COPY go.mod go.sum ./
RUN go mod download && go mod verify

# ============================================
# Stage 2: Builder
# ============================================
FROM deps AS builder

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -o /app/server ./main.go

# ============================================
# Stage 3: Production (Distroless)
# ============================================
FROM gcr.io/distroless/static-debian12:nonroot AS production

WORKDIR /app

COPY --from=builder /app/server        /app/server
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /app/migrations    /app/migrations

USER nonroot:nonroot

EXPOSE 50059
EXPOSE 2112

CMD ["/app/server"]

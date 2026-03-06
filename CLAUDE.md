# Ticket Microservices Platform

## Architecture

Microservice architecture: **Next.js frontend** → **Node.js gateway** → **gRPC** → **Go/Java backend services** → **PostgreSQL**

### Services & Ports

| Service | Language | gRPC Port | Role |
|---------|----------|-----------|------|
| gateway | Node.js (Express) | — | API gateway, REST → gRPC translation |
| frontend | Next.js 15 (App Router) | — | Web UI (port 3000) |
| auth-service | Go | 50051 | Authentication, RBAC |
| user-service | Go | 50052 | User profiles, addresses |
| event-service | Go | 50053 | Events, zones, pricing, seats, schedules |
| ticket-service | Go | 50054 | Ticket types, availability, reservation |
| booking-service | Go | 50058 | Bookings, cancellation |
| checkin-service | Go | 50059 | Event check-in |
| invoice-service | Go | 50060 | Invoice generation |
| payment-service | Go | 50062 | Payment processing |
| realtime-service | Go | — | WebSocket (port 50070) |
| booking-worker | Go | — | Background job processor (Redis queue) |
| email-worker | Go | — | Email sending (Kafka consumer) |

### Infrastructure

- **PostgreSQL**: Auth DB (port 50432), Main DB (port 50433)
- **Redis**: Cache (50379), Queue (50380), PubSub (50381)
- **Kafka**: Message broker (50092)

## Data Flow

```
Frontend (Next.js) → Gateway (REST/JSON) → gRPC → Backend Service → PostgreSQL
                                ↕
                    Case transform: snake_case ↔ camelCase
                    (gateway auto-converts both directions)
```

**Important:** Frontend uses camelCase, backend uses snake_case. The gateway middleware handles conversion automatically. Do NOT manually convert cases.

## Development Commands

```bash
# Frontend
cd frontend && npm run dev          # Dev server (port 3000)
cd frontend && npx tsc --noEmit     # Type check

# Go services (any service)
cd event-service && go build ./...  # Build check
cd event-service && air             # Hot reload dev server

# Gateway
cd gateway && npm run dev           # Dev server (port 53000)

# Docker (databases + infra)
docker compose -f deploy/docker-compose.dev.yml up -d
```

## Proto Workflow

Proto files live in `shared-lib/protos/`. After modifying a `.proto` file:

```bash
# Generate Go code (use temp dir to handle go_package path)
TMPDIR=$(mktemp -d) && \
protoc --go_out=$TMPDIR --go-grpc_out=$TMPDIR \
  --proto_path=shared-lib/protos <proto_file>.proto && \
cp $TMPDIR/github.com/booking-system/shared-lib/protos/<pkg>/*.pb.go \
  <service>/internal/protos/<pkg>/
```

**Rules:**
- Never reuse proto field numbers (even after deleting a field)
- Always increment field numbers for new fields
- Regenerate `.pb.go` after every proto change

## DB Migrations

Location: `<service>/database/migrations/`
Naming: `{NNN}_{description}.{up|down}.sql` (e.g., `005_add_event_status.up.sql`)

**Rules:**
- Always create both `.up.sql` and `.down.sql`
- Use `IF NOT EXISTS` / `IF EXISTS` for idempotency
- After adding columns, update: Go model → repository queries → gRPC controller

## Frontend Conventions

### Atomic Design (component hierarchy)

```
components/ui/          → Atoms (shadcn/ui — never modify directly)
components/molecules/   → Combine 2-3 atoms, pure presentation, no data fetching
components/organisms/   → Complex sections, may have state, grouped by domain
  ├── layout/           → Navbar, Sidebar, Footer
  ├── auth/             → Login/Register forms
  ├── events/           → Public event components
  ├── org-events/       → Org event management components
  ├── admin-*/          → Admin dashboard components
  └── shared/           → DataTable, generic organisms
app/(public)/           → Public pages (no auth)
app/(auth)/             → Login, Register
app/(dashboard)/        → User dashboard (individual role)
app/(org)/              → Organization pages
app/(admin)/            → Admin pages
```

### Key Frontend Paths

```
lib/api/client.ts       → Axios instance with auth interceptor
lib/api/endpoints.ts    → All API endpoint URLs
lib/api/types/          → TypeScript interfaces per domain
lib/api/queries/        → TanStack Query hooks per domain
lib/validators/         → Zod schemas
stores/                 → Zustand stores
```

### Stack

- **UI**: shadcn/ui (Radix primitives) + Tailwind CSS + `cn()` from `lib/utils`
- **Forms**: react-hook-form + zod + zodResolver
- **Server state**: TanStack Query v5
- **Client state**: Zustand
- **Icons**: lucide-react
- **Toasts**: sonner

## Common Gotchas

1. **Hybrid ID pattern**: DB uses `BIGSERIAL` (internal) + `UUID` (public-facing `public_id`). Frontend always receives/sends the UUID.
2. **Gateway injects user context**: `req.user.id` is added by auth middleware — used as `created_by` / `organization_id` in create requests.
3. **Event status**: Events are created as `"draft"`, explicitly published via `/events/:id/publish`.
4. **Proto `go_package`**: Points to `github.com/booking-system/shared-lib/protos/<pkg>` — causes nested output dirs. Use temp dir + copy pattern.

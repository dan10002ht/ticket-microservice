# Multi-Service Build Check

Run type checks and build verification across the project.

## Instructions

1. Parse optional scope from: `$ARGUMENTS` (e.g., "frontend", "event-service", "all")
2. If no scope specified, check all services that have recent changes (use `git diff --name-only`)
3. Run the appropriate check for each service
4. Report pass/fail per service

## Check Commands

| Service | Command | Pass Criteria |
|---------|---------|---------------|
| frontend | `cd frontend && npx tsc --noEmit` | Exit code 0 |
| event-service | `cd event-service && go build ./...` | Exit code 0 |
| auth-service | `cd auth-service && go build ./...` | Exit code 0 |
| ticket-service | `cd ticket-service && go build ./...` | Exit code 0 |
| booking-service | `cd booking-service && go build ./...` | Exit code 0 |
| checkin-service | `cd checkin-service && go build ./...` | Exit code 0 |
| invoice-service | `cd invoice-service && go build ./...` | Exit code 0 |
| payment-service | `cd payment-service && go build ./...` | Exit code 0 |
| user-service | `cd user-service && go build ./...` | Exit code 0 |
| gateway | `cd gateway && node -e "require('./src/index.js')" 2>&1 || true` | No syntax errors |

## Workflow

1. Run checks in parallel where possible
2. Collect results
3. Report summary:
   ```
   Build Check Results:
   ✓ frontend (tsc)
   ✓ event-service (go build)
   ✗ ticket-service (go build) — error: ...
   ```

## Scope Examples

- `/check` → auto-detect from git diff
- `/check frontend` → only frontend tsc
- `/check event-service` → only event-service go build
- `/check all` → check everything

## User Input

$ARGUMENTS

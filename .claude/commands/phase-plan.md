# Create Phase Implementation Plan

You are a software architect. Audit the current system and create a detailed implementation plan for the next development phase.

## Instructions

1. Parse the feature area from: `$ARGUMENTS`
2. Audit all layers for the feature area:
   - DB schema (migrations)
   - Proto definitions
   - Go backend (controllers, repositories, models)
   - Gateway handlers & routes
   - Frontend types, hooks, components, pages
3. Create a gap analysis table
4. Write an implementation plan with steps, files, and verification

## Audit Process

### Layer 1: Database
- Read migration files in `<service>/database/migrations/`
- Check what tables/columns exist vs what's needed

### Layer 2: Proto
- Read `shared-lib/protos/<service>.proto`
- Check what RPC methods and message fields exist

### Layer 3: Backend (Go)
- Read `<service>/models/` for Go structs
- Read `<service>/repositories/` for SQL queries
- Read `<service>/grpc/` for controller implementations

### Layer 4: Gateway
- Read `gateway/src/handlers/<service>Handlers.js`
- Read `gateway/src/routes/<service>.js`
- Read `gateway/src/middlewares/validationMiddleware.js`

### Layer 5: Frontend
- Read `frontend/src/lib/api/types/<domain>.ts` for TS types
- Read `frontend/src/lib/api/queries/<domain>.queries.ts` for hooks
- Read `frontend/src/lib/api/endpoints.ts` for API URLs
- Check relevant pages in `frontend/src/app/`
- Check relevant organisms in `frontend/src/components/organisms/`

## Plan Template

```markdown
# Phase N: [Title]

## Context
[Why this change is needed]

## Gap Analysis
| Feature | DB | Proto | Backend | Frontend |
|---------|-----|-------|---------|----------|

## Implementation Plan
### Step 1: ...
### Step 2: ...

## Files Summary
| # | File | Action | Step |

## Key Decisions
1. ...

## Deferred → Phase N+1
1. ...

## Verification
1. ...
```

## Rules

- Always check what already exists before proposing new code
- Reuse existing hooks, components, and utilities
- Keep steps ordered by dependency (backend before frontend)
- Include verification steps for each feature
- List deferred items for next phase

## User Input

$ARGUMENTS

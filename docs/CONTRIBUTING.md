# Contributing Guidelines

## Code Style

### Node.js (auth-service, gateway)

- Functions over classes
- camelCase file names
- ES Modules (`import/export`)
- Use `yarn` as package manager

### Go (user, event, ticket, realtime, booking-worker, email-worker)

- Standard Go conventions
- snake_case file names
- Go 1.22+

### Java (booking-service, payment-service)

- Java 17, Spring Boot 3.2
- PascalCase classes
- Maven build tool
- Spring profiles for environment config (`application-dev.yml`)

---

## Branch Naming

```
feature/<ticket-id>-short-description
fix/<ticket-id>-short-description
refactor/<description>
docs/<description>
```

---

## Commit Messages

Use conventional commits:

```
feat: add ticket reservation endpoint
fix: resolve payment webhook timeout
refactor: extract booking saga steps
docs: update event service API docs
```

---

## Pull Request Workflow

1. Create a feature branch from `main`
2. Make changes with clear, atomic commits
3. Ensure all tests pass locally
4. Update documentation if adding/changing APIs
5. Submit PR with description of changes and test plan
6. Get review approval before merging

---

## Documentation Standards

- All service documentation lives in `docs/services/<service-name>/`
- Architecture decisions go in `docs/architecture/`
- Setup guides go in `docs/guides/`
- Update the service's README when adding new features or changing APIs
- Keep port references in sync with `scripts/lib/service-config.sh`

---

## Testing

- **Node.js**: Jest
- **Go**: Standard `go test`
- **Java**: JUnit 5 + Spring Boot Test

Run service-specific tests before submitting PRs:

```bash
# Node.js services
cd auth-service && yarn test

# Go services
cd user-service && go test ./...

# Java services
cd booking-service && mvn test
```

---

## Adding a New Service

1. Create service directory at project root: `<service-name>/`
2. Add service to `scripts/lib/service-config.sh` (ports, health check config)
3. Add startup logic to `scripts/dev-all.sh`
4. Create docs at `docs/services/<service-name>/README.md`
5. If exposing REST API through gateway:
   - Add route file in `gateway/src/routes/`
   - Add handler in `gateway/src/handlers/`
   - Add swagger docs in `gateway/src/swagger/`
   - Mount in `gateway/src/services/routeService.js`
6. Update `docs/services/README.md` with new service entry

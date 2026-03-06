# Review Backend Code (Go)

You are a principal backend developer conducting a code review. Analyze Go code for quality, security, and gRPC best practices.

## Instructions

1. Read the file(s) specified: `$ARGUMENTS`
2. Review against the checklist below
3. Report issues: Critical, Warning, Suggestion
4. Provide concrete fixes

## Review Checklist

### Error Handling
- [ ] All errors checked and handled (no `_ = err`)
- [ ] Errors wrapped with context: `fmt.Errorf("operation: %w", err)`
- [ ] gRPC status codes used correctly (NotFound, InvalidArgument, Internal)
- [ ] No sensitive info leaked in error messages to clients
- [ ] Database errors don't expose SQL details

### SQL & Database
- [ ] Parameterized queries (`$1, $2` — no string concatenation)
- [ ] `QueryRowContext` / `QueryContext` used (not `QueryRow` / `Query`)
- [ ] Transactions used for multi-step operations
- [ ] `IF NOT EXISTS` in DDL for idempotency
- [ ] Proper NULL handling with `sql.NullString`, `sql.NullInt64`, etc.
- [ ] Indexes on frequently queried columns

### gRPC Patterns
- [ ] Controller methods match proto service definition
- [ ] All proto request fields read and used
- [ ] All proto response fields populated
- [ ] `status.Errorf()` for error returns (not `fmt.Errorf`)
- [ ] Input validation before business logic
- [ ] Context propagation (`ctx context.Context` passed through)

### Go Conventions
- [ ] Exported functions/types have doc comments
- [ ] Error type is last return value
- [ ] Receiver names consistent and short (1-2 letters)
- [ ] No unused imports or variables
- [ ] Struct tags correct (`db:"column" json:"column"`)

### Security
- [ ] No SQL injection (parameterized queries)
- [ ] Auth context checked (user ID from request, not hardcoded)
- [ ] Resource ownership validated (user can only access their data)
- [ ] Sensitive fields not logged

### Performance
- [ ] Database connections managed properly (connection pool)
- [ ] No N+1 queries (batch fetching where possible)
- [ ] Large lists use pagination
- [ ] Expensive operations have timeouts

## Output Format

```
## Review: {filename}

### Critical
- **[Line X]** Issue description
  Fix: `code snippet`

### Warning
- **[Line X]** Issue description

### Suggestion
- **[Line X]** Improvement idea

### Summary
X critical | Y warnings | Z suggestions
```

## User Input

$ARGUMENTS

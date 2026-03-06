---
paths:
  - "gateway/src/**/*.js"
---

# Gateway Rules

## Architecture
The gateway is a Node.js (Express) app that translates REST → gRPC calls.

## Request Flow
```
Client (JSON) → Express route → Auth middleware → Validation → Handler → gRPC client → Backend service
                                                                ↕
                                                    Case transform middleware
                                                    (snake_case ↔ camelCase)
```

## Handler Pattern
```javascript
const handler = async (req, res) => {
    const result = await grpcClients.serviceName.methodName({
        ...req.body,              // or req.params, req.query
        created_by: req.user.id,  // inject auth context
    });
    sendSuccessResponse(res, 200, result, req.correlationId);
};
```

## Key Files
- `src/config/index.js` — Service URLs, ports, env vars
- `src/grpc/clients.js` — gRPC client connections
- `src/handlers/<service>Handlers.js` — Request handlers per service
- `src/routes/<service>.js` — Express route definitions
- `src/middlewares/validationMiddleware.js` — express-validator rules
- `src/middlewares/authMiddleware.js` — JWT verification

## Case Transform
- Gateway auto-converts: request body camelCase → snake_case, response snake_case → camelCase
- Do NOT manually convert cases in handlers
- `req.user` is populated by auth middleware: `{ id, email, role }`

## Adding a New Route
1. Add handler function in `src/handlers/<service>Handlers.js`
2. Add route in `src/routes/<service>.js` with appropriate middleware
3. Add validation in `src/middlewares/validationMiddleware.js` if needed

# 🚀 API Checklist: Event Creation Flow (Organization)

## 1. Organization Dashboard API

- [ ] **GET /api/organizations/:orgId/dashboard**
  - [ ] Route: `/src/routes/organization.js`
  - [ ] Handler: `getOrganizationDashboardHandler`
  - [ ] gRPC: `organizationService.getDashboard`
  - [ ] Validation: orgId param
  - [ ] Swagger docs
  - [ ] Unit tests
  - [ ] Response: Dashboard stats + recent events

## 2. Event Creation APIs

- [ ] **POST /api/events**
  - [x] Route: `/src/routes/event.js`
  - [x] Handler: `createEventHandler`
  - [x] gRPC: `eventService.createEvent`
  - [x] Validation: request body
  - [ ] Swagger docs
  - [ ] Unit tests
  - [x] Response: Created event with ID
  - [x] Authorization: requireRole(['organization'])

- [ ] **PUT /api/events/:eventId/draft**
  - [x] Route: `/src/routes/event.js`
  - [x] Handler: `saveEventDraftHandler`
  - [ ] gRPC: `eventService.updateEvent`
  - [x] Validation: eventId + partial data
  - [ ] Swagger docs
  - [ ] Unit tests
  - [x] Response: Updated draft event

## 3. Layout Management APIs

- [ ] **POST /api/events/:eventId/layout**
  - [ ] Route: `/src/routes/event.js`
  - [ ] Handler: `saveEventLayoutHandler`
  - [ ] gRPC: `eventService.updateEventLayout`
  - [ ] Validation: eventId + layout data
  - [ ] Swagger docs
  - [ ] Unit tests
  - [ ] Response: Updated layout

## 4. Pricing Management APIs

- [ ] **POST /api/events/:eventId/pricing**
  - [ ] Route: `/src/routes/event.js`
  - [ ] Handler: `saveEventPricingHandler`
  - [ ] gRPC: `pricingService.setEventPricing`
  - [ ] Validation: eventId + pricing data
  - [ ] Swagger docs
  - [ ] Unit tests
  - [ ] Response: Updated pricing

## 5. Event Publishing API

- [ ] **POST /api/events/:eventId/publish**
  - [ ] Route: `/src/routes/event.js`
  - [ ] Handler: `publishEventHandler`
  - [ ] gRPC: `eventService.publishEvent`
  - [ ] Validation: eventId + publish options
  - [ ] Swagger docs
  - [ ] Unit tests
  - [ ] Response: Published event

## 6. Event Templates API

- [ ] **GET /api/events/templates**
  - [ ] Route: `/src/routes/event.js`
  - [ ] Handler: `getEventTemplatesHandler`
  - [ ] gRPC: `eventService.getTemplates`
  - [ ] Validation: none
  - [ ] Swagger docs
  - [ ] Unit tests
  - [ ] Response: Array of templates

## 7. Event Duplication API

- [ ] **POST /api/events/:eventId/duplicate**
  - [ ] Route: `/src/routes/event.js`
  - [ ] Handler: `duplicateEventHandler`
  - [ ] gRPC: `eventService.duplicateEvent`
  - [ ] Validation: eventId + options
  - [ ] Swagger docs
  - [ ] Unit tests
  - [ ] Response: New duplicated event

---

## 🔧 General Implementation Checklist

- [ ] Express route setup
- [ ] Authentication middleware
- [ ] Validation middleware
- [ ] Error handling middleware
- [ ] gRPC client setup (event, pricing, organization)
- [ ] Swagger/OpenAPI documentation
- [ ] Unit & integration tests
- [ ] Response formatting
- [ ] Pagination (nếu cần)
- [ ] Error response formatting

---

**Status**: 🟡 Planning
**Priority**: HIGH
**Dependencies**: Organization service, Event service, Pricing service

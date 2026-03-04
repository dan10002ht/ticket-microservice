export const API_ENDPOINTS = {
  auth: {
    login: "/auth/login",
    register: "/auth/register",
    registerEmail: "/auth/register/email",
    registerOAuth: "/auth/register/oauth",
    refresh: "/auth/refresh",
    logout: "/auth/logout",
    validate: "/auth/validate",
    oauthLogin: "/auth/oauth/login",
    sendVerificationEmail: "/auth/send-verification-email",
    verifyUser: "/auth/verify-user",
    resendVerificationEmail: "/auth/resend-verification-email",
    verifyEmailToken: "/auth/verify-email-token",
    forgotPassword: "/auth/forgot-password",
    resetPassword: "/auth/reset-password",
    permissions: "/auth/permissions",
    permissionsCheck: "/auth/permissions/check",
    permissionsResource: "/auth/permissions/resource",
    permissionsBatch: "/auth/permissions/batch",
    roles: "/auth/roles",
    user: (userId: string) => `/auth/users/${userId}`,
    health: "/auth/health",
  },

  events: {
    list: "/events",
    detail: (id: string) => `/events/${id}`,
    draft: (id: string) => `/events/${id}/draft`,
    publish: (id: string) => `/events/${id}/publish`,
    duplicate: (id: string) => `/events/${id}/duplicate`,
    templates: "/events/templates",
    // Zones
    zones: (eventId: string) => `/events/${eventId}/zones`,
    zone: (eventId: string, zoneId: string) =>
      `/events/${eventId}/zones/${zoneId}`,
    // Seats
    seats: (eventId: string) => `/events/${eventId}/seats`,
    seatsBulk: (eventId: string) => `/events/${eventId}/seats/bulk`,
    seat: (eventId: string, seatId: string) =>
      `/events/${eventId}/seats/${seatId}`,
    // Pricing
    pricing: (eventId: string) => `/events/${eventId}/pricing`,
    pricingDetail: (eventId: string, pricingId: string) =>
      `/events/${eventId}/pricing/${pricingId}`,
    pricingByZone: (eventId: string, zoneId: string) =>
      `/events/${eventId}/pricing/zone/${zoneId}`,
    pricingCalculate: (eventId: string) =>
      `/events/${eventId}/pricing/calculate`,
    pricingDiscount: (eventId: string) =>
      `/events/${eventId}/pricing/discount`,
    // Availability
    availability: (eventId: string) => `/events/${eventId}/availability`,
    availabilityZone: (eventId: string, zoneId: string) =>
      `/events/${eventId}/availability/zones/${zoneId}`,
    availabilitySeat: (eventId: string, seatId: string) =>
      `/events/${eventId}/availability/seats/${seatId}`,
    availabilityBlock: (eventId: string) =>
      `/events/${eventId}/availability/block`,
    availabilityRelease: (eventId: string) =>
      `/events/${eventId}/availability/release`,
  },

  bookings: {
    list: "/bookings",
    detail: (id: string) => `/bookings/${id}`,
    cancel: (id: string) => `/bookings/${id}/cancel`,
    confirm: (id: string) => `/bookings/${id}/confirm`,
    seatsReserve: "/bookings/seats/reserve",
    seatsRelease: "/bookings/seats/release",
    adminList: "/bookings/admin/list",
  },

  tickets: {
    list: "/tickets",
    detail: (id: string) => `/tickets/${id}`,
    types: (eventId: string) => `/tickets/types/${eventId}`,
    typesCreate: "/tickets/types",
    typeDetail: (typeId: string) => `/tickets/types/${typeId}`,
    availability: (eventId: string) => `/tickets/availability/${eventId}`,
    available: (eventId: string) => `/tickets/available/${eventId}`,
    reserve: (eventId: string) => `/tickets/reserve/${eventId}`,
    release: "/tickets/release",
  },

  payments: {
    list: "/payments",
    create: "/payments",
    detail: (id: string) => `/payments/${id}`,
    capture: (id: string) => `/payments/${id}/capture`,
    cancel: (id: string) => `/payments/${id}/cancel`,
    refund: (id: string) => `/payments/${id}/refund`,
    refunds: (id: string) => `/payments/${id}/refunds`,
    refundDetail: (refundId: string) => `/payments/refunds/${refundId}`,
    methods: "/payments/methods",
    adminList: "/payments/admin/list",
  },

  users: {
    profile: "/users/profile",
    addresses: "/users/addresses",
    address: (id: string) => `/users/addresses/${id}`,
    adminList: "/users/admin/list",
    adminDetail: (userId: string) => `/users/admin/${userId}`,
    adminCreate: "/users/admin",
  },

  health: {
    check: "/health",
    ready: "/health/ready",
    live: "/health/live",
  },
} as const;

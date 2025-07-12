import { Counter } from 'prom-client';

export const authAttempts = new Counter({
  name: 'auth_attempts_total',
  help: 'Total authentication attempts',
  labelNames: ['method', 'status', 'user_type'],
});

import { AsyncLocalStorage } from 'async_hooks';

const storage = new AsyncLocalStorage();

/**
 * Request-scoped context storage.
 * Allows correlationId to flow automatically into gRPC calls
 * without threading it through every function argument.
 */
export const requestContext = {
  /**
   * Run a function within a context bound to the given store.
   * Call this once per incoming HTTP request.
   */
  run: (ctx, fn) => storage.run(ctx, fn),

  /** Return the context for the current async execution chain, or null. */
  get: () => storage.getStore() ?? null,
};

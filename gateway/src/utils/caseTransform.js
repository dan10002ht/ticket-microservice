/**
 * Deep transform object keys between snake_case and camelCase.
 * Uses LRU-style key cache to avoid repeated regex on hot paths.
 */

const camelCache = Object.create(null);
const snakeCache = Object.create(null);

function snakeToCamel(str) {
  if (str in camelCache) return camelCache[str];
  // Fast path: no underscore → already camelCase
  if (str.indexOf('_') === -1) {
    camelCache[str] = str;
    return str;
  }
  const result = str.replace(/_([a-z])/g, (_, c) => c.toUpperCase());
  camelCache[str] = result;
  return result;
}

function camelToSnake(str) {
  if (str in snakeCache) return snakeCache[str];
  const result = str.replace(/[A-Z]/g, (c) => `_${c.toLowerCase()}`);
  snakeCache[str] = result;
  return result;
}

function transformKeys(obj, keyFn) {
  if (obj === null || obj === undefined || typeof obj !== 'object') {
    return obj;
  }

  if (Array.isArray(obj)) {
    const len = obj.length;
    const out = new Array(len);
    for (let i = 0; i < len; i++) {
      out[i] = transformKeys(obj[i], keyFn);
    }
    return out;
  }

  // Skip Date, Buffer, etc.
  const proto = Object.getPrototypeOf(obj);
  if (proto !== Object.prototype && proto !== null) {
    return obj;
  }

  const keys = Object.keys(obj);
  const out = {};
  for (let i = 0, len = keys.length; i < len; i++) {
    const k = keys[i];
    out[keyFn(k)] = transformKeys(obj[k], keyFn);
  }
  return out;
}

/**
 * Deep transform all keys from snake_case to camelCase
 * Used for: gRPC response -> REST API response
 */
export function toCamelCase(obj) {
  return transformKeys(obj, snakeToCamel);
}

/**
 * Deep transform all keys from camelCase to snake_case
 * Used for: REST API request -> gRPC request
 */
export function toSnakeCase(obj) {
  return transformKeys(obj, camelToSnake);
}

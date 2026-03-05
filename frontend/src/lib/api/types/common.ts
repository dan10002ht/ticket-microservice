/** Generic API response wrapper */
export interface ApiResponse<T> {
  data: T;
  meta?: {
    correlationId?: string;
    timestamp?: string;
  };
}

/** Normalized paginated response (backend formats vary — normalize in query hooks) */
export interface PaginatedResponse<T> {
  items: T[];
  total: number;
  page: number;
  limit: number;
  totalPages?: number;
}

/** Pagination request params */
export interface PaginationParams {
  page?: number;
  limit?: number;
}

/** Typed API error matching gateway error shape */
export interface ApiError {
  error: {
    code: string;
    message: string;
    details?: unknown;
  };
  meta?: {
    correlationId?: string;
    timestamp?: string;
  };
  /** HTTP status code (set by error interceptor) */
  statusCode?: number;
}

/** Sort direction */
export type SortDirection = "asc" | "desc";

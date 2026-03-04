/** Generic API response wrapper */
export interface ApiResponse<T> {
  success?: boolean;
  data?: T;
  message?: string;
  correlationId?: string;
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
  error: string;
  correlationId?: string;
  code?: string;
  statusCode?: number;
  timestamp?: string;
  details?: unknown;
}

/** Sort direction */
export type SortDirection = "asc" | "desc";

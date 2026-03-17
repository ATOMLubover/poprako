/**
 * 通用分页查询参数。
 */
export interface PaginationQuery {
  offset: number;
  limit: number;
}

/**
 * Swagger 中常见的 include 数组查询参数。
 */
export interface IncludeQuery {
  includes?: string[];
}

/**
 * 统一 API 错误结构。
 */
export interface ApiErrorPayload {
  code?: string;
  message: string;
}

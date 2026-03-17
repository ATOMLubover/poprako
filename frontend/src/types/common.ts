/**
 * 文件用途：定义全局通用类型，例如分页参数与统一错误结构。
 */
/**
 * 通用分页查询参数。
 */
export interface PaginationQuery {
  /**
   * 分页偏移量，表示从第几条记录开始读取。
   */
  offset: number;
  /**
   * 单页条数，表示本次请求最多返回多少条记录。
   */
  limit: number;
}

/**
 * Swagger 中常见的 include 数组查询参数。
 */
export interface IncludeQuery {
  /**
   * 关联展开字段列表，对应 swagger 中 includes[] 查询参数。
   */
  includes?: string[];
}

/**
 * 统一 API 错误结构。
 */
export interface ApiErrorPayload {
  /**
   * 业务错误码，由后端定义；缺省时通常表示仅返回了通用错误。
   */
  code?: string;
  /**
   * 面向前端展示或日志记录的错误消息。
   */
  message: string;
}

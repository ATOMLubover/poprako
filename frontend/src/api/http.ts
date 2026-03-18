/**
 * 文件用途：Axios 统一请求层，集中处理鉴权、错误处理与 HTTP 方法封装。
 */
import axios, {
  AxiosError,
  AxiosHeaders,
  type AxiosInstance,
  type AxiosRequestConfig,
  type AxiosResponse,
  type InternalAxiosRequestConfig,
} from "axios";
import type { ApiErrorPayload } from "../types/common";

interface ApiResponseEnvelope<T> {
  code: number;
  message: string;
  data?: T;
}

function resolveApiBaseURL(): string {
  const configuredBaseURL = import.meta.env.VITE_API_BASE_URL?.trim();
  if (configuredBaseURL) {
    return configuredBaseURL;
  }

  return "/api/v1";
}

/**
 * 统一请求客户端。
 * 该类负责请求头注入、错误处理、响应解包与通用 HTTP 方法封装。
 */
class ApiHttpClient {
  /**
   * Axios 实例对象，负责执行所有 HTTP 请求。
   */
  private readonly instance: AxiosInstance;

  /**
   * 构造请求客户端并初始化拦截器。
   */
  public constructor() {
    this.instance = axios.create({
      baseURL: resolveApiBaseURL(),
      timeout: 15_000,
    });
    this.setupInterceptors();
  }

  /**
   * 安装请求与响应拦截器。
   */
  private setupInterceptors(): void {
    this.instance.interceptors.request.use((config) =>
      this.handleRequest(config),
    );
    this.instance.interceptors.response.use(
      (response) => response,
      (error: AxiosError<ApiErrorPayload>) => this.handleResponseError(error),
    );
  }

  /**
   * 处理请求拦截逻辑：自动挂载访问令牌。
   */
  private handleRequest(
    config: InternalAxiosRequestConfig,
  ): InternalAxiosRequestConfig {
    const token = localStorage.getItem("access_token");
    if (token) {
      if (!(config.headers instanceof AxiosHeaders)) {
        config.headers = AxiosHeaders.from(config.headers);
      }
      config.headers.set("Authorization", `Bearer ${token}`);
    }
    return config;
  }

  /**
   * 处理响应错误：标准化错误信息并处理未授权场景。
   */
  private handleResponseError(
    error: AxiosError<ApiErrorPayload>,
  ): Promise<never> {
    const statusCode = error.response?.status;
    const message =
      error.response?.data?.message ?? error.message ?? "请求失败";

    if (statusCode === 401) {
      localStorage.removeItem("access_token");
      if (location.pathname !== "/login") {
        location.href = "/login";
      }
    }

    return Promise.reject(new Error(message));
  }

  /**
   * 统一请求入口。
   */
  public async request<T>(config: AxiosRequestConfig): Promise<T> {
    const response: AxiosResponse<ApiResponseEnvelope<T>> =
      await this.instance.request<ApiResponseEnvelope<T>>(config);
    const responseBody = response.data;

    if (responseBody.code !== 200) {
      return Promise.reject(new Error(responseBody.message || "请求失败"));
    }

    return (responseBody.data ?? (undefined as T)) as T;
  }

  /**
   * 发送 GET 请求。
   */
  public async get<T>(url: string, params?: object): Promise<T> {
    return this.request<T>({
      url,
      method: "GET",
      params,
      paramsSerializer: (queryParams) => this.serializeQuery(queryParams),
    });
  }

  /**
   * 发送 POST 请求。
   */
  public async post<T, B = unknown>(url: string, body?: B): Promise<T> {
    return this.request<T>({
      url,
      method: "POST",
      data: body,
    });
  }

  /**
   * 发送 PUT 请求。
   */
  public async put<T, B = unknown>(url: string, body?: B): Promise<T> {
    return this.request<T>({
      url,
      method: "PUT",
      data: body,
    });
  }

  /**
   * 发送 PATCH 请求。
   */
  public async patch<T, B = unknown>(url: string, body?: B): Promise<T> {
    return this.request<T>({
      url,
      method: "PATCH",
      data: body,
    });
  }

  /**
   * 发送 DELETE 请求。
   */
  public async delete<T>(url: string): Promise<T> {
    return this.request<T>({
      url,
      method: "DELETE",
    });
  }

  /**
   * 序列化查询参数，兼容 swagger 中 includes[] 的数组格式。
   */
  private serializeQuery(queryParams: object): string {
    const searchParams = new URLSearchParams();
    Object.entries(queryParams as Record<string, unknown>).forEach(
      ([key, value]) => {
        if (value === undefined || value === null) {
          return;
        }
        if (Array.isArray(value)) {
          value.forEach((singleValue) =>
            searchParams.append(`${key}[]`, String(singleValue)),
          );
          return;
        }
        searchParams.append(key, String(value));
      },
    );
    return searchParams.toString();
  }
}

/**
 * 导出全局单例，作为项目唯一请求出口。
 */
export const httpClient = new ApiHttpClient();

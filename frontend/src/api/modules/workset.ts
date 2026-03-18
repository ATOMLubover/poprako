/**
 * 文件用途：封装工作集（workset）相关接口。
 */
import { httpClient } from "../http";
import type { PaginationQuery } from "../../types/common";
import type { WorksetInfo } from "../../types/domain";

/**
 * 工作集查询参数，对应 GET /worksets。
 */
export interface WorksetListQuery extends PaginationQuery {
  /** 所属团队 ID。 */
  team_id: string;
}

/**
 * 获取工作集列表请求类型。
 */
export type GetWorksetListRequest = WorksetListQuery;

/**
 * 创建工作集参数，对应 swagger 的 value.CreateWorksetArgs。
 */
export interface CreateWorksetArgs {
  /** 所属团队 ID。 */
  team_id: string;
  /** 工作集名称。 */
  name: string;
  /** 工作集描述。 */
  description?: string;
}

/**
 * 创建工作集请求类型。
 */
export type CreateWorksetRequest = CreateWorksetArgs;

/**
 * 获取工作集列表响应类型。
 */
export type GetWorksetListResponse = WorksetInfo[];

/**
 * 创建工作集响应类型。
 */
export type CreateWorksetResponse = WorksetInfo;

/**
 * 获取工作集列表，对应 GET /worksets。
 * 请求类型：GetWorksetListRequest。
 * 返回类型：GetWorksetListResponse。
 */
export async function getWorksetList(
  worksetListQuery: GetWorksetListRequest,
): Promise<GetWorksetListResponse> {
  return httpClient.get<GetWorksetListResponse>("/worksets", worksetListQuery);
}

/**
 * 创建工作集，对应 POST /worksets。
 * 请求类型：CreateWorksetRequest。
 * 返回类型：CreateWorksetResponse。
 */
export async function createWorkset(
  createWorksetArgs: CreateWorksetRequest,
): Promise<CreateWorksetResponse> {
  return httpClient.post<CreateWorksetResponse, CreateWorksetRequest>(
    "/worksets",
    createWorksetArgs,
  );
}

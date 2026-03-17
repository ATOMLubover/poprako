/**
 * 文件用途：封装团队（team）相关接口。
 */
import { httpClient } from "../http";
import type { PaginationQuery } from "../../types/common";
import type { TeamInfo } from "../../types/domain";

/**
 * 团队列表查询参数。
 */
export type TeamListQuery = PaginationQuery;

/**
 * 获取团队列表的请求类型。
 */
export type GetTeamListRequest = TeamListQuery;

/**
 * 创建团队参数，对应 swagger 的 value.CreateTeamArgs。
 */
export interface CreateTeamArgs {
  /** 团队名称。 */
  name: string;
  /** 团队简介。 */
  description?: string;
}

/**
 * 创建团队请求类型。
 */
export type CreateTeamRequest = CreateTeamArgs;

/**
 * 获取当前用户团队列表响应类型。
 */
export type GetMyTeamsResponse = TeamInfo[];

/**
 * 获取团队列表响应类型。
 */
export type GetTeamListResponse = TeamInfo[];

/**
 * 创建团队响应类型。
 */
export type CreateTeamResponse = TeamInfo;

/**
 * 获取当前用户团队列表，对应 GET /teams/mine。
 * 请求类型：GetTeamListRequest。
 * 返回类型：GetMyTeamsResponse。
 */
export async function getMyTeams(
  teamListQuery: GetTeamListRequest = {
    offset: 0,
    limit: 20,
  },
): Promise<GetMyTeamsResponse> {
  return httpClient.get<GetMyTeamsResponse>("/teams/mine", teamListQuery);
}

/**
 * 获取全部团队列表，对应 GET /teams。
 * 请求类型：GetTeamListRequest。
 * 返回类型：GetTeamListResponse。
 */
export async function getTeamList(
  teamListQuery: GetTeamListRequest,
): Promise<GetTeamListResponse> {
  return httpClient.get<GetTeamListResponse>("/teams", teamListQuery);
}

/**
 * 创建团队，对应 POST /teams。
 * 请求类型：CreateTeamRequest。
 * 返回类型：CreateTeamResponse。
 */
export async function createTeam(
  createTeamArgs: CreateTeamRequest,
): Promise<CreateTeamResponse> {
  return httpClient.post<CreateTeamResponse, CreateTeamRequest>(
    "/teams",
    createTeamArgs,
  );
}

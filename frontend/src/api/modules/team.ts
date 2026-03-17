import { httpClient } from '@/api/http';
import type { PaginationQuery } from '@/types/common';
import type { TeamInfo } from '@/types/domain';

/**
 * 团队列表查询参数。
 */
export type TeamListQuery = PaginationQuery;

/**
 * 创建团队参数，对应 swagger 的 value.CreateTeamArgs。
 */
export interface CreateTeamArgs {
  name: string;
  description?: string;
}

/**
 * 获取当前用户团队列表，对应 GET /teams/mine。
 */
export async function getMyTeams(): Promise<TeamInfo[]> {
  return httpClient.get<TeamInfo[]>('/teams/mine');
}

/**
 * 获取全部团队列表，对应 GET /teams。
 */
export async function getTeamList(teamListQuery: TeamListQuery): Promise<TeamInfo[]> {
  return httpClient.get<TeamInfo[]>('/teams', teamListQuery);
}

/**
 * 创建团队，对应 POST /teams。
 */
export async function createTeam(createTeamArgs: CreateTeamArgs): Promise<TeamInfo> {
  return httpClient.post<TeamInfo, CreateTeamArgs>('/teams', createTeamArgs);
}

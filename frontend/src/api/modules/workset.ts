import { httpClient } from '@/api/http';
import type { PaginationQuery } from '@/types/common';
import type { WorksetInfo } from '@/types/domain';

/**
 * 工作集查询参数，对应 GET /worksets。
 */
export interface WorksetListQuery extends PaginationQuery {
  team_id: string;
}

/**
 * 创建工作集参数，对应 swagger 的 value.CreateWorksetArgs。
 */
export interface CreateWorksetArgs {
  team_id: string;
  name: string;
  description?: string;
}

/**
 * 获取工作集列表，对应 GET /worksets。
 */
export async function getWorksetList(worksetListQuery: WorksetListQuery): Promise<WorksetInfo[]> {
  return httpClient.get<WorksetInfo[]>('/worksets', worksetListQuery);
}

/**
 * 创建工作集，对应 POST /worksets。
 */
export async function createWorkset(createWorksetArgs: CreateWorksetArgs): Promise<WorksetInfo> {
  return httpClient.post<WorksetInfo, CreateWorksetArgs>('/worksets', createWorksetArgs);
}

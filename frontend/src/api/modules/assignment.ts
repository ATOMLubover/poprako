import { httpClient } from '@/api/http';
import type { IncludeQuery, PaginationQuery } from '@/types/common';
import type { AssignmentInfo } from '@/types/domain';

/**
 * 分配列表查询参数，对应 GET /assignments。
 */
export interface AssignmentListQuery extends PaginationQuery, IncludeQuery {
  chapter_id: string;
}

/**
 * 创建分配参数，对应 swagger 的 value.CreateChapterAssignmentArgs。
 */
export interface CreateAssignmentArgs {
  chapter_id: string;
  user_id: string;
  role: string;
}

/**
 * 获取章节分配列表，对应 GET /assignments。
 */
export async function getAssignmentList(assignmentListQuery: AssignmentListQuery): Promise<AssignmentInfo[]> {
  return httpClient.get<AssignmentInfo[]>('/assignments', assignmentListQuery);
}

/**
 * 获取我的分配列表，对应 GET /assignments/mine。
 */
export async function getMyAssignments(paginationQuery: PaginationQuery): Promise<AssignmentInfo[]> {
  return httpClient.get<AssignmentInfo[]>('/assignments/mine', paginationQuery);
}

/**
 * 创建分配记录，对应 POST /assignments。
 */
export async function createAssignment(createAssignmentArgs: CreateAssignmentArgs): Promise<AssignmentInfo> {
  return httpClient.post<AssignmentInfo, CreateAssignmentArgs>('/assignments', createAssignmentArgs);
}

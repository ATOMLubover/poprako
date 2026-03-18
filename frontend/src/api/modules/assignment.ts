/**
 * 文件用途：封装章节分配（assignment）相关接口。
 */
import { httpClient } from "../http";
import type { IncludeQuery, PaginationQuery } from "../../types/common";
import type { AssignmentInfo } from "../../types/domain";

/**
 * 分配列表查询参数，对应 GET /assignments。
 */
export interface AssignmentListQuery extends PaginationQuery, IncludeQuery {
  /** 所属章节 ID。 */
  chapter_id: string;
}

/**
 * 获取分配列表请求类型。
 */
export type GetAssignmentListRequest = AssignmentListQuery;

/**
 * 获取我的分配请求类型。
 */
export type GetMyAssignmentsRequest = PaginationQuery;

/**
 * 创建分配参数，对应 swagger 的 value.CreateChapterAssignmentArgs。
 */
export interface CreateAssignmentArgs {
  /** 所属章节 ID。 */
  chapter_id: string;
  /** 被分配用户 ID。 */
  user_id: string;
  /** 分配角色，例如 translator / reviewer。 */
  role: string;
}

/**
 * 创建分配请求类型。
 */
export type CreateAssignmentRequest = CreateAssignmentArgs;

/**
 * 获取分配列表响应类型。
 */
export type GetAssignmentListResponse = AssignmentInfo[];

/**
 * 获取我的分配响应类型。
 */
export type GetMyAssignmentsResponse = AssignmentInfo[];

/**
 * 创建分配响应类型。
 */
export type CreateAssignmentResponse = AssignmentInfo;

/**
 * 获取章节分配列表，对应 GET /assignments。
 * 请求类型：GetAssignmentListRequest。
 * 返回类型：GetAssignmentListResponse。
 */
export async function getAssignmentList(
  assignmentListQuery: GetAssignmentListRequest,
): Promise<GetAssignmentListResponse> {
  return httpClient.get<GetAssignmentListResponse>(
    "/assignments",
    assignmentListQuery,
  );
}

/**
 * 获取我的分配列表，对应 GET /assignments/mine。
 * 请求类型：GetMyAssignmentsRequest。
 * 返回类型：GetMyAssignmentsResponse。
 */
export async function getMyAssignments(
  paginationQuery: GetMyAssignmentsRequest,
): Promise<GetMyAssignmentsResponse> {
  const assignmentList = await httpClient.get<GetMyAssignmentsResponse>(
    "/assignments/mine",
    paginationQuery,
  );

  return Array.isArray(assignmentList) ? assignmentList : [];
}

/**
 * 创建分配记录，对应 POST /assignments。
 * 请求类型：CreateAssignmentRequest。
 * 返回类型：CreateAssignmentResponse。
 */
export async function createAssignment(
  createAssignmentArgs: CreateAssignmentRequest,
): Promise<CreateAssignmentResponse> {
  return httpClient.post<CreateAssignmentResponse, CreateAssignmentRequest>(
    "/assignments",
    createAssignmentArgs,
  );
}

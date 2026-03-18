/**
 * 文件用途：封装章节（chapter）相关接口。
 */
import { httpClient } from "../http";
import type { IncludeQuery, PaginationQuery } from "../../types/common";
import type { ChapterInfo } from "../../types/domain";

/**
 * 章节列表查询参数，对应 GET /chapters。
 */
export interface ChapterListQuery extends PaginationQuery, IncludeQuery {
  /** 所属漫画 ID。 */
  comic_id: string;
}

/**
 * 获取章节列表请求类型。
 */
export type GetChapterListRequest = ChapterListQuery;

/**
 * 创建章节参数，对应 swagger 的 value.CreateChapterArgs。
 */
export interface CreateChapterArgs {
  /** 所属漫画 ID。 */
  comic_id: string;
  /** 章节标题。 */
  title: string;
  /** 章节序号。 */
  index: number;
}

/**
 * 创建章节请求类型。
 */
export type CreateChapterRequest = CreateChapterArgs;

/**
 * 获取章节列表响应类型。
 */
export type GetChapterListResponse = ChapterInfo[];

/**
 * 创建章节响应类型。
 */
export type CreateChapterResponse = ChapterInfo;

/**
 * 获取章节列表，对应 GET /chapters。
 * 请求类型：GetChapterListRequest。
 * 返回类型：GetChapterListResponse。
 */
export async function getChapterList(
  chapterListQuery: GetChapterListRequest,
): Promise<GetChapterListResponse> {
  return httpClient.get<GetChapterListResponse>("/chapters", chapterListQuery);
}

/**
 * 创建章节，对应 POST /chapters。
 * 请求类型：CreateChapterRequest。
 * 返回类型：CreateChapterResponse。
 */
export async function createChapter(
  createChapterArgs: CreateChapterRequest,
): Promise<CreateChapterResponse> {
  return httpClient.post<CreateChapterResponse, CreateChapterRequest>(
    "/chapters",
    createChapterArgs,
  );
}

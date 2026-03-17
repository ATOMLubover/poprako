/**
 * 文件用途：封装漫画（comic）相关接口。
 */
import { httpClient } from "../http";
import type { PaginationQuery } from "../../types/common";
import type { ComicInfo } from "../../types/domain";

/**
 * 漫画列表查询参数，对应 GET /comics。
 */
export interface ComicListQuery extends PaginationQuery {
  /** 所属工作集 ID。 */
  workset_id: string;
}

/**
 * 获取漫画列表请求类型。
 */
export type GetComicListRequest = ComicListQuery;

/**
 * 创建漫画参数，对应 swagger 的 value.CreateComicArgs。
 */
export interface CreateComicArgs {
  /** 所属工作集 ID。 */
  workset_id: string;
  /** 漫画标题。 */
  title: string;
}

/**
 * 创建漫画请求类型。
 */
export type CreateComicRequest = CreateComicArgs;

/**
 * 获取漫画列表响应类型。
 */
export type GetComicListResponse = ComicInfo[];

/**
 * 创建漫画响应类型。
 */
export type CreateComicResponse = ComicInfo;

/**
 * 获取漫画列表，对应 GET /comics。
 * 请求类型：GetComicListRequest。
 * 返回类型：GetComicListResponse。
 */
export async function getComicList(
  comicListQuery: GetComicListRequest,
): Promise<GetComicListResponse> {
  return httpClient.get<GetComicListResponse>("/comics", comicListQuery);
}

/**
 * 创建漫画，对应 POST /comics。
 * 请求类型：CreateComicRequest。
 * 返回类型：CreateComicResponse。
 */
export async function createComic(
  createComicArgs: CreateComicRequest,
): Promise<CreateComicResponse> {
  return httpClient.post<CreateComicResponse, CreateComicRequest>(
    "/comics",
    createComicArgs,
  );
}

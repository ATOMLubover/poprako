import { httpClient } from '@/api/http';
import type { PaginationQuery } from '@/types/common';
import type { ComicInfo } from '@/types/domain';

/**
 * 漫画列表查询参数，对应 GET /comics。
 */
export interface ComicListQuery extends PaginationQuery {
  workset_id: string;
}

/**
 * 创建漫画参数，对应 swagger 的 value.CreateComicArgs。
 */
export interface CreateComicArgs {
  workset_id: string;
  title: string;
}

/**
 * 获取漫画列表，对应 GET /comics。
 */
export async function getComicList(comicListQuery: ComicListQuery): Promise<ComicInfo[]> {
  return httpClient.get<ComicInfo[]>('/comics', comicListQuery);
}

/**
 * 创建漫画，对应 POST /comics。
 */
export async function createComic(createComicArgs: CreateComicArgs): Promise<ComicInfo> {
  return httpClient.post<ComicInfo, CreateComicArgs>('/comics', createComicArgs);
}

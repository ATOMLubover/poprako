import { httpClient } from '@/api/http';
import type { IncludeQuery, PaginationQuery } from '@/types/common';
import type { ChapterInfo } from '@/types/domain';

/**
 * 章节列表查询参数，对应 GET /chapters。
 */
export interface ChapterListQuery extends PaginationQuery, IncludeQuery {
  comic_id: string;
}

/**
 * 创建章节参数，对应 swagger 的 value.CreateChapterArgs。
 */
export interface CreateChapterArgs {
  comic_id: string;
  title: string;
  index: number;
}

/**
 * 获取章节列表，对应 GET /chapters。
 */
export async function getChapterList(chapterListQuery: ChapterListQuery): Promise<ChapterInfo[]> {
  return httpClient.get<ChapterInfo[]>('/chapters', chapterListQuery);
}

/**
 * 创建章节，对应 POST /chapters。
 */
export async function createChapter(createChapterArgs: CreateChapterArgs): Promise<ChapterInfo> {
  return httpClient.post<ChapterInfo, CreateChapterArgs>('/chapters', createChapterArgs);
}

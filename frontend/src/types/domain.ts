/**
 * 用户信息，对应 swagger 的 value.UserInfo。
 */
export interface UserInfo {
  id: string;
  username: string;
  qq: string;
  avatar?: string;
}

/**
 * 团队信息，对应 swagger 的 value.TeamInfo。
 */
export interface TeamInfo {
  id: string;
  name: string;
  avatar?: string;
  description?: string;
  created_at?: string;
}

/**
 * 工作集信息，对应 swagger 的 value.WorksetInfo。
 */
export interface WorksetInfo {
  id: string;
  team_id: string;
  name: string;
  description?: string;
}

/**
 * 漫画信息，对应 swagger 的 value.ComicInfo。
 */
export interface ComicInfo {
  id: string;
  workset_id: string;
  title: string;
  cover?: string;
}

/**
 * 章节信息，对应 swagger 的 value.ChapterInfo。
 */
export interface ChapterInfo {
  id: string;
  comic_id: string;
  title: string;
  index: number;
}

/**
 * 分配信息，对应 swagger 的 value.AssignmentInfo。
 */
export interface AssignmentInfo {
  id: string;
  chapter_id: string;
  role: string;
  user_id: string;
}

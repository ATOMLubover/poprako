/**
 * 文件用途：定义与 Swagger 领域对象对应的核心业务类型。
 */
/**
 * 用户信息，对应 swagger 的 value.UserInfo。
 */
export interface UserInfo {
  /** 用户唯一标识。 */
  id: string;
  /** 用户昵称或显示名。 */
  username: string;
  /** 用户绑定的 QQ 号码。 */
  qq: string;
  /** 用户头像 URL，未设置时可能为空。 */
  avatar?: string;
}

/**
 * 团队信息，对应 swagger 的 value.TeamInfo。
 */
export interface TeamInfo {
  /** 汉化组唯一标识。 */
  id: string;
  /** 汉化组名称。 */
  name: string;
  /** 汉化组头像 URL，未设置时可能为空。 */
  avatar?: string;
  /** 汉化组描述信息。 */
  description?: string;
  /** 汉化组创建时间，通常为 ISO 时间字符串。 */
  created_at?: string;
}

/**
 * 工作集信息，对应 swagger 的 value.WorksetInfo。
 */
export interface WorksetInfo {
  /** 工作集唯一标识。 */
  id: string;
  /** 所属汉化组 ID。 */
  team_id: string;
  /** 工作集名称。 */
  name: string;
  /** 工作集描述信息。 */
  description?: string;
}

/**
 * 漫画信息，对应 swagger 的 value.ComicInfo。
 */
export interface ComicInfo {
  /** 漫画唯一标识。 */
  id: string;
  /** 所属工作集 ID。 */
  workset_id: string;
  /** 漫画标题。 */
  title: string;
  /** 漫画封面 URL，未上传时可能为空。 */
  cover?: string;
}

/**
 * 章节信息，对应 swagger 的 value.ChapterInfo。
 */
export interface ChapterInfo {
  /** 章节唯一标识。 */
  id: string;
  /** 所属漫画 ID。 */
  comic_id: string;
  /** 章节标题。 */
  title: string;
  /** 章节序号。 */
  index: number;
}

/**
 * 分配信息，对应 swagger 的 value.AssignmentInfo。
 */
export interface AssignmentInfo {
  /** 分配记录唯一标识。 */
  id: string;
  /** 所属章节 ID。 */
  chapter_id: string;
  /** 分配角色，例如 translator、reviewer、done。 */
  role: string;
  /** 被分配用户 ID。 */
  user_id: string;
}

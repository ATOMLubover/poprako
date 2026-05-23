// ─── Auth ────────────────────────────────────────────────────────────────────

export interface LoginUserRes {
  user_id: string;
  token: string;
}

export interface RegUserRes {
  user_id: string;
  token: string;
}

// ─── Team ────────────────────────────────────────────────────────────────────

export interface CreateTeamRes {
  id: string;
}

export interface TeamInfo {
  id: string;
  name: string;
  description: string;
}

// ─── Team announcement / comment ─────────────────────────────────────────────

export interface EmbeddedUserInfo {
  id: string;
  qid: string;
  nickname: string;
}

export interface AnnouncementVal {
  id: string;
  team_id: string;
  user_id: string;
  user: EmbeddedUserInfo | null;
  title: string;
  content: string;
  created_at: number;
}

export interface AnnouncementCreatedRes {
  id: string;
}

export interface CommentVal {
  id: string;
  team_id: string;
  user_id: string;
  user: EmbeddedUserInfo | null;
  content: string;
  created_at: number;
}

export interface CommentCreatedRes {
  id: string;
}

// ─── Member ──────────────────────────────────────────────────────────────────

export interface CreateMemberRes {
  id: string;
}

// ─── Invitation ──────────────────────────────────────────────────────────────

export interface InvitationInfo {
  id: string;
  team_id: string;
  invitee_qid: string;
  invitation_code: string;
  roles: number;
  pending: boolean;
  created_at: number;
}

// ─── Workset ─────────────────────────────────────────────────────────────────

export interface CreateWorksetRes {
  id: string;
}

// ─── Comic ───────────────────────────────────────────────────────────────────

export interface CreateComicRes {
  id: string;
}

// ─── Chapter ─────────────────────────────────────────────────────────────────

export interface CreateChapterRes {
  id: string;
}

// ─── Page ────────────────────────────────────────────────────────────────────

export interface PageCreationResult {
  page_id: string;
  put_url: string;
}

export interface ReserveChapterPagesRes {
  creations: PageCreationResult[];
}

// ─── Assignment ──────────────────────────────────────────────────────────────

export interface CreateAssignmentRes {
  id: string;
}

// ─── Unit ────────────────────────────────────────────────────────────────────

export interface UnitInfo {
  id: string;
  page_id: string;
  index: number;
  x_coord: number;
  y_coord: number;
  is_bubble: boolean;
  translated_text: string | null;
  translator_id: string | null;
  translator_comment: string | null;
  is_proofread: boolean;
  proofread_text: string | null;
  proofreader_id: string | null;
  proofreader_comment: string | null;
}

// ─── Generic API envelope ────────────────────────────────────────────────────

export interface ApiResponse<T> {
  code: number;
  message: string;
  data: T;
}

// ─── Session state (accumulated across phases) ────────────────────────────────

export interface SeedState {
  adminToken: string;
  adminUserID: string;

  teamID: string;

  translatorToken: string;
  translatorUserID: string;

  proofreaderToken: string;
  proofreaderUserID: string;

  worksetID: string;
  comicID: string;
  chapterID: string;
  pageIDs: string[];
}

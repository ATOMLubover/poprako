/**
 * Phase 9 — Team board
 *
 * 1. Admin creates one team announcement
 * 2. Translator creates one team comment
 * 3. Member lists announcements and verifies embedded user
 * 4. Member lists comments and verifies embedded user
 */

import { api, apiGet, logStep, logOk, logInfo } from "../client";
import type {
  AnnouncementCreatedRes,
  AnnouncementVal,
  CommentCreatedRes,
  CommentVal,
  SeedState,
} from "../types";

function pickAnnouncementById(
  items: AnnouncementVal[],
  id: string,
): AnnouncementVal {
  for (const item of items) {
    if (item.id === id) {
      return item;
    }
  }

  throw new Error(`announcement not found: ${id}`);
}

function pickCommentById(items: CommentVal[], id: string): CommentVal {
  for (const item of items) {
    if (item.id === id) {
      return item;
    }
  }

  throw new Error(`comment not found: ${id}`);
}

export async function phase9TeamBoard(state: SeedState): Promise<void> {
  logStep(
    "Phase 9",
    "Team board — announcement/comment create + list + embedded user",
  );

  // ── 9.1 Admin creates one announcement ────────────────────────────────────
  const announcementRes = await api<AnnouncementCreatedRes>(
    "POST",
    "/announcements",
    {
      team_id: state.teamID,
      title: "Seed Announcement",
      content: "Welcome to the seed workspace",
    },
    { token: state.adminToken },
  );

  logOk("Announcement created", { announcement_id: announcementRes.id });

  // ── 9.2 Translator creates one comment ────────────────────────────────────
  const commentRes = await api<CommentCreatedRes>(
    "POST",
    "/comments",
    {
      team_id: state.teamID,
      content: "Hello team, translator online",
    },
    { token: state.translatorToken },
  );

  logOk("Comment created", { comment_id: commentRes.id });

  // ── 9.3 Proofreader lists announcements and verifies embedded user ────────
  const announcementList = await apiGet<AnnouncementVal[]>("/announcements", {
    token: state.proofreaderToken,
    query: {
      team_id: state.teamID,
      offset: 0,
      limit: 3,
    },
  });

  const createdAnnouncement = pickAnnouncementById(
    announcementList,
    announcementRes.id,
  );

  if (!createdAnnouncement.user) {
    throw new Error("announcement list user should not be null");
  }

  if (createdAnnouncement.user.id !== state.adminUserID) {
    throw new Error("announcement user id mismatch");
  }

  logOk("Announcement list includes embedded user", {
    user_id: createdAnnouncement.user.id,
  });

  // ── 9.4 Admin lists comments and verifies embedded user ───────────────────
  const commentList = await apiGet<CommentVal[]>("/comments", {
    token: state.adminToken,
    query: {
      team_id: state.teamID,
      offset: 0,
      limit: 20,
    },
  });

  const createdComment = pickCommentById(commentList, commentRes.id);

  if (!createdComment.user) {
    throw new Error("comment list user should not be null");
  }

  if (createdComment.user.id !== state.translatorUserID) {
    throw new Error("comment user id mismatch");
  }

  logOk("Comment list includes embedded user", {
    user_id: createdComment.user.id,
  });

  logInfo("Announcement/comment full-chain checks passed");
}

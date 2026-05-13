/**
 * Phase 2 — Create content
 *
 * workset → comic → chapter → pages (reserved + marked uploaded)
 *
 * The backend synchronously assigns the chapter creator as REVIEWER during
 * chapter creation (ChapterCreatorAssignedEvent, PubTypeSync), so pages can
 * be reserved immediately after the chapter is created.
 * Unit writes are still blocked until Phase 4 advances the workflow.
 */

import { api, logStep, logOk, logInfo } from "../client";
import { ROLE } from "../config";
import type {
  CreateWorksetRes,
  CreateComicRes,
  CreateChapterRes,
  ReserveChapterPagesRes,
  SeedState,
} from "../types";

const PAGE_COUNT = 3;

export async function phase2Content(state: SeedState): Promise<void> {
  logStep("Phase 2", "Create content — workset → comic → chapter → pages");

  // ── 2.1 Create workset ────────────────────────────────────────────────────
  const worksetRes = await api<CreateWorksetRes>(
    "POST",
    "/worksets",
    {
      team_id: state.teamID,
      name: "Seed Workset",
      description: "Integration test workset",
    },
    { token: state.adminToken },
  );

  state.worksetID = worksetRes.id;
  logOk("Workset created", { workset_id: worksetRes.id });

  // ── 2.2 Create comic ──────────────────────────────────────────────────────
  const comicRes = await api<CreateComicRes>(
    "POST",
    "/comics",
    {
      workset_id: state.worksetID,
      title: "Seed Comic",
      author: "Seed Author",
      description: "Integration test comic",
    },
    { token: state.adminToken },
  );

  state.comicID = comicRes.id;
  logOk("Comic created", { comic_id: comicRes.id });

  // ── 2.3 Create chapter ────────────────────────────────────────────────────
  const chapterRes = await api<CreateChapterRes>(
    "POST",
    "/chapters",
    {
      comic_id: state.comicID,
      subtitle: "Chapter 1",
    },
    { token: state.adminToken },
  );

  state.chapterID = chapterRes.id;
  logOk("Chapter created", { chapter_id: chapterRes.id });

  // The creator gets REVIEWER automatically, but page upload confirmation needs
  // a chapter-level RAW_PROVIDER assignment as well.
  await api<null>(
    "PUT",
    "/assignments",
    {
      chapter_id: state.chapterID,
      user_id: state.adminUserID,
      role_mask: ROLE.RAW_PROVIDER | ROLE.REVIEWER,
    },
    { token: state.adminToken },
  );

  logOk("Admin chapter assignment updated with RAW_PROVIDER");

  // ── 2.4 Reserve pages ─────────────────────────────────────────────────────
  // The REVIEWER assignment for the chapter creator was made synchronously by
  // the backend during chapter creation, so this call succeeds immediately.
  const pagesRes = await api<ReserveChapterPagesRes>(
    "POST",
    "/pages/reserve",
    {
      page_count: PAGE_COUNT,
      file_extension: "png",
    },
    { token: state.adminToken, query: { chapter_id: state.chapterID } },
  );

  state.pageIDs = pagesRes.creations.map((c) => c.page_id);
  logOk(`${PAGE_COUNT} pages reserved`, { page_ids: state.pageIDs });

  // ── 2.5 Mark each page as uploaded ───────────────────────────────────────
  for (const pageID of state.pageIDs) {
    await api<null>(
      "POST",
      `/pages/${pageID}/image/uploaded`,
      null,
      { token: state.adminToken },
    );
  }

  logOk(`All ${PAGE_COUNT} pages marked as uploaded`);
  logInfo("Note: workflow is still Pending — no unit writes allowed yet");
}

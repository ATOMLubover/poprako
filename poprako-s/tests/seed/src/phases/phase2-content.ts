/**
 * Phase 2 — Create content
 *
 * workset → comic → chapter → pages (reserved + marked uploaded)
 *
 * Note: At this point NO ONE can write units yet.
 *       Assignments are needed first (Phase 3).
 */

import { api, logStep, logOk, logInfo } from "../client";
import type {
  CreateWorksetRes,
  CreateComicRes,
  CreateChapterRes,
  ReserveChapterPagesRes,
  SeedState,
} from "../types";

const PAGE_COUNT = 3;
const RESERVE_RETRY_TIMES = 60;
const RESERVE_RETRY_DELAY_MS = 500;

async function reservePagesWithRetry(
  chapterID: string,
  token: string,
): Promise<ReserveChapterPagesRes> {
  let lastErr: unknown;

  for (let i = 1; i <= RESERVE_RETRY_TIMES; i++) {
    try {
      return await api<ReserveChapterPagesRes>(
        "POST",
        "/pages",
        {
          chapter_id: chapterID,
          page_count: PAGE_COUNT,
          extension: "png",
        },
        { token },
      );
    } catch (err) {
      lastErr = err;
      if (i < RESERVE_RETRY_TIMES) {
        logInfo("Reserve pages not ready yet, retrying", {
          attempt: i,
          max_attempts: RESERVE_RETRY_TIMES,
        });
        await Bun.sleep(RESERVE_RETRY_DELAY_MS);
      }
    }
  }

  throw new Error(
    "Reserve pages failed after retries. " +
      "If chapter creator auto-reviewer assignment is asynchronous, it may not be applied yet; " +
      "or the server build still does not include that behavior. " +
      `Last error: ${String(lastErr)}`,
  );
}

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

  // ── 2.4 Reserve pages ─────────────────────────────────────────────────────
  const pagesRes = await reservePagesWithRetry(
    state.chapterID,
    state.adminToken,
  );

  state.pageIDs = pagesRes.creations.map((c) => c.page_id);
  logOk(`${PAGE_COUNT} pages reserved`, { page_ids: state.pageIDs });

  // ── 2.5 Mark each page as uploaded ───────────────────────────────────────
  for (const pageID of state.pageIDs) {
    await api<null>(
      "PUT",
      `/pages/${pageID}`,
      { id: pageID, is_uploaded: true },
      { token: state.adminToken },
    );
  }

  logOk(`All ${PAGE_COUNT} pages marked as uploaded`);
  logInfo("Note: workflow is still Pending — no unit writes allowed yet");
}

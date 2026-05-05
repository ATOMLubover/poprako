/**
 * Phase 3 — Create Assignments
 *
 * Admin (who has REVIEWER role) creates assignments for:
 *   - translator → TRANSLATOR role
 *   - proofreader → PROOFREADER role
 *
 * ⚠️ Assignments MUST be created BEFORE any unit writes.
 */

import { api, logStep, logOk, logInfo } from "../client";
import { ROLE } from "../config";
import type { SeedState } from "../types";

export async function phase3Assignments(state: SeedState): Promise<void> {
  logStep("Phase 3", "Create assignments — translator + proofreader");

  // ── 3.1 Assign translator ─────────────────────────────────────────────────
  await api<null>(
    "PUT",
    "/assignments",
    {
      chapter_id: state.chapterID,
      user_id: state.translatorUserID,
      role_mask: ROLE.TRANSLATOR,
    },
    { token: state.adminToken },
  );

  logOk("Translator assigned");

  // ── 3.2 Assign proofreader ────────────────────────────────────────────────
  await api<null>(
    "PUT",
    "/assignments",
    {
      chapter_id: state.chapterID,
      user_id: state.proofreaderUserID,
      role_mask: ROLE.PROOFREADER,
    },
    { token: state.adminToken },
  );

  logOk("Proofreader assigned");

  logInfo("Note: assignments exist but workflow gates still block unit writes");
}

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
import type { CreateAssignmentRes, SeedState } from "../types";

export async function phase3Assignments(state: SeedState): Promise<void> {
  logStep("Phase 3", "Create assignments — translator + proofreader");

  // ── 3.1 Assign translator ─────────────────────────────────────────────────
  const assignTranslator = await api<CreateAssignmentRes>(
    "POST",
    "/assignments",
    {
      chapter_id: state.chapterID,
      user_id: state.translatorUserID,
      roles: ROLE.TRANSLATOR,
    },
    { token: state.adminToken },
  );

  logOk("Translator assigned", { assignment_id: assignTranslator.id });

  // ── 3.2 Assign proofreader ────────────────────────────────────────────────
  const assignProofreader = await api<CreateAssignmentRes>(
    "POST",
    "/assignments",
    {
      chapter_id: state.chapterID,
      user_id: state.proofreaderUserID,
      roles: ROLE.PROOFREADER,
    },
    { token: state.adminToken },
  );

  logOk("Proofreader assigned", { assignment_id: assignProofreader.id });

  logInfo("Note: assignments exist but workflow gates still block unit writes");
}

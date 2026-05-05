/**
 * Phase 8 — Finalize workflow
 *
 * proofread_complete → typeset_start → typeset_complete
 *   → review_complete → publish_complete
 *
 * All transitions are driven by admin (reviewer role).
 */

import { api, logStep, logOk } from "../client";
import { WORKFLOW } from "../config";
import type { SeedState } from "../types";

const TRANSITIONS = [
  WORKFLOW.PROOFREAD_COMPLETE,
  WORKFLOW.TYPESET_START,
  WORKFLOW.TYPESET_COMPLETE,
  WORKFLOW.REVIEW_COMPLETE,
  WORKFLOW.PUBLISH_COMPLETE,
] as const;

export async function phase8Finalize(state: SeedState): Promise<void> {
  logStep("Phase 8", "Finalize workflow — proofread → publish");

  for (const transition of TRANSITIONS) {
    await api<null>(
      "PUT",
      `/chapters/${state.chapterID}`,
      { workflow_transition: transition },
      { token: state.adminToken },
    );

    logOk("Workflow transition fired", { transition });
  }

  logOk("Chapter fully published ✓");
}

/**
 * Phase 6 — Advance workflow to proofread stage
 *
 * translate_complete → closes translator window
 * proofread_start    → opens proofreader window
 */

import { api, logStep, logOk } from "../client";
import { WORKFLOW } from "../config";
import type { SeedState } from "../types";

export async function phase6WorkflowProofread(state: SeedState): Promise<void> {
  logStep("Phase 6", "Workflow advance — translate_complete → proofread_start");

  // ── 6.1 translate_complete ────────────────────────────────────────────────
  await api<null>(
    "PATCH",
    `/chapters/${state.chapterID}`,
    { workflow_transition: WORKFLOW.TRANSLATE_COMPLETE },
    { token: state.adminToken },
  );

  logOk("Workflow transition fired", {
    transition: WORKFLOW.TRANSLATE_COMPLETE,
  });

  // ── 6.2 proofread_start ───────────────────────────────────────────────────
  await api<null>(
    "PATCH",
    `/chapters/${state.chapterID}`,
    { workflow_transition: WORKFLOW.PROOFREAD_START },
    { token: state.adminToken },
  );

  logOk("Workflow transition fired", { transition: WORKFLOW.PROOFREAD_START });
  logOk("Proofreader can now write units");
}

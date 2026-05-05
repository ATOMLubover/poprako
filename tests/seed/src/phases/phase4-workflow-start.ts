/**
 * Phase 4 — Advance workflow to translation stage
 *
 * upload_complete → unlocks upload phase
 * translate_start → unlocks translator unit writes
 */

import { api, logStep, logOk } from "../client";
import { WORKFLOW } from "../config";
import type { SeedState } from "../types";

export async function phase4WorkflowStart(state: SeedState): Promise<void> {
  logStep("Phase 4", "Workflow advance — upload_complete → translate_start");

  // ── 4.1 upload_complete ───────────────────────────────────────────────────
  await api<null>(
    "PUT",
    `/chapters/${state.chapterID}`,
    { workflow_transition: WORKFLOW.UPLOAD_COMPLETE },
    { token: state.adminToken },
  );

  logOk("Workflow transition fired", { transition: WORKFLOW.UPLOAD_COMPLETE });

  // ── 4.2 translate_start ───────────────────────────────────────────────────
  await api<null>(
    "PUT",
    `/chapters/${state.chapterID}`,
    { workflow_transition: WORKFLOW.TRANSLATE_START },
    { token: state.adminToken },
  );

  logOk("Workflow transition fired", { transition: WORKFLOW.TRANSLATE_START });
  logOk("Translator can now write units");
}

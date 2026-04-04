/**
 * Phase 7 — Proofreader reviews units
 *
 * 1. Re-fetch units as proofreader (mandatory — cannot reuse translator's fetch)
 * 2. Patch first unit: set proofread_text + is_proofread = true
 * 3. Patch second unit
 */

import { api, apiGet, logStep, logOk, logInfo } from "../client";
import type { UnitInfo, SeedState } from "../types";

export async function phase7Proofread(state: SeedState): Promise<void> {
  logStep("Phase 7", "Proofreader reviews units on page 0");

  const pageID = state.pageIDs[0];

  // ── 7.1 Re-fetch units as proofreader ─────────────────────────────────────
  // Must use proofreader's token and must re-fetch (not reuse phase-5 result).
  const units = await apiGet<UnitInfo[]>("/units", {
    token: state.proofreaderToken,
    query: { page_id: pageID },
  });

  logOk("Units re-fetched as proofreader", { count: units?.length ?? 0 });

  if (!units || units.length === 0) {
    throw new Error("Phase 7: server returned no units for proofreader");
  }

  logInfo(
    "Units to proofread",
    units.map((u) => ({ id: u.id, index: u.index })),
  );

  // ── 7.2 Patch all units ───────────────────────────────────────────────────
  const patches = units.map((u) => ({
    id: u.id,
    proofread_text: `${u.translated_text ?? ""} [已校对]`,
    is_proofread: true,
    proofreader_comment: "校对通过",
  }));

  await api<null>(
    "PUT",
    "/units",
    {
      page_id: pageID,
      unit_diff: { patch: patches },
    },
    { token: state.proofreaderToken },
  );

  logOk(`Patched ${patches.length} units as proofread`);
}

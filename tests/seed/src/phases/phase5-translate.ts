/**
 * Phase 5 — Translator writes units
 *
 * For the first page:
 *   1. Insert 2 units (client-provided IDs, server stores them)
 *   2. Re-fetch units to get server-authoritative IDs
 *   3. Patch one unit to add a translator comment
 *
 * ⚠️ Always re-fetch before patching — never assume stored ID equals
 *    what was sent (round-trip is mandatory per seed contract).
 */

import { api, apiGet, logStep, logOk, logInfo } from "../client";
import type { UnitInfo, SeedState } from "../types";

export async function phase5Translate(state: SeedState): Promise<void> {
  logStep("Phase 5", "Translator writes units on page 0");

  const pageID = state.pageIDs[0];

  // ── 5.1 Insert units ──────────────────────────────────────────────────────
  // Client provides IDs; server stores them as-is.
  // We still must re-fetch afterwards for authoritative data.
  const insertID1 = crypto.randomUUID();
  const insertID2 = crypto.randomUUID();

  await api<null>(
    "POST",
    `/pages/${pageID}/units`,
    {
      page_id: pageID,
      difference: {
        page_id: pageID,
        operations: [
          {
            local_id: insertID1,
            is_bubble: true,
            is_proofread: false,
            x_coord: 100,
            y_coord: 200,
            translated_text: "你好，世界！",
          },
          {
            local_id: insertID2,
            is_bubble: true,
            is_proofread: false,
            x_coord: 150,
            y_coord: 300,
            translated_text: "这是第二个气泡框。",
          },
        ],
        candidate_order: [insertID1, insertID2],
      },
    },
    { token: state.translatorToken },
  );

  logOk("Units inserted (2 bubbles)");

  // ── 5.2 Re-fetch units (mandatory round-trip) ─────────────────────────────
  const unitList = await apiGet<{ units: UnitInfo[] }>(`/pages/${pageID}/units`, {
    token: state.translatorToken,
  });
  const units = unitList.units;

  logOk(`Units fetched from server`, { count: units?.length ?? 0 });

  if (!units || units.length === 0) {
    throw new Error("Phase 5: server returned no units after insert");
  }

  const firstUnit = units[0]!;
  logInfo("First unit", {
    id: firstUnit.id,
    translated_text: firstUnit.translated_text,
  });

  // ── 5.3 Patch first unit — add translator comment ─────────────────────────
  await api<null>(
    "POST",
    `/pages/${pageID}/units`,
    {
      page_id: pageID,
      difference: {
        page_id: pageID,
        operations: [
          {
            id: firstUnit.id,
            is_bubble: firstUnit.is_bubble,
            is_proofread: firstUnit.is_proofread,
            x_coord: firstUnit.x_coord,
            y_coord: firstUnit.y_coord,
            translated_text: firstUnit.translated_text,
            translator_comment: "译者注：此处使用了双关语。",
          },
        ],
        candidate_order: units.map((u) => u.id),
      },
    },
    { token: state.translatorToken },
  );

  logOk("Unit patched with translator comment", { unit_id: firstUnit.id });
}

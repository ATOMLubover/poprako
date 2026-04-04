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
    "PUT",
    "/units",
    {
      page_id: pageID,
      unit_diff: {
        insert: [
          {
            id: insertID1,
            index: 0,
            x_coord: 100,
            y_coord: 200,
            is_bubble: true,
            translated_text: "你好，世界！",
          },
          {
            id: insertID2,
            index: 1,
            x_coord: 150,
            y_coord: 300,
            is_bubble: true,
            translated_text: "这是第二个气泡框。",
          },
        ],
      },
    },
    { token: state.translatorToken },
  );

  logOk("Units inserted (2 bubbles)");

  // ── 5.2 Re-fetch units (mandatory round-trip) ─────────────────────────────
  const units = await apiGet<UnitInfo[]>("/units", {
    token: state.translatorToken,
    query: { page_id: pageID },
  });

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
    "PUT",
    "/units",
    {
      page_id: pageID,
      unit_diff: {
        patch: [
          {
            id: firstUnit.id,
            translator_comment: "译者注：此处使用了双关语。",
          },
        ],
      },
    },
    { token: state.translatorToken },
  );

  logOk("Unit patched with translator comment", { unit_id: firstUnit.id });
}

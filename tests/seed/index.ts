/**
 * Seed — Happy-path integration test runner
 *
 * Runs all phases in order, accumulating state in a single `SeedState` object.
 * Each phase logs clearly so failures are easy to pinpoint.
 *
 * Usage:
 *   bun run seed
 *   bun run index.ts
 */

import { phase0Bootstrap } from "./src/phases/phase0-bootstrap";
import { phase1Users } from "./src/phases/phase1-users";
import { phase2Content } from "./src/phases/phase2-content";
import { phase3Assignments } from "./src/phases/phase3-assignments";
import { phase4WorkflowStart } from "./src/phases/phase4-workflow-start";
import { phase5Translate } from "./src/phases/phase5-translate";
import { phase6WorkflowProofread } from "./src/phases/phase6-workflow-proofread";
import { phase7Proofread } from "./src/phases/phase7-proofread";
import { phase8Finalize } from "./src/phases/phase8-finalize";
import { phase9TeamBoard } from "./src/phases/phase9-team-board";
import type { SeedState } from "./src/types";

const RESET = "\x1b[0m";
const BOLD = "\x1b[1m";
const GREEN = "\x1b[32m";
const RED = "\x1b[31m";
const CYAN = "\x1b[36m";

async function run(): Promise<void> {
  console.log(
    `\n${BOLD}${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${RESET}`,
  );
  console.log(
    `${BOLD}${CYAN}  PopRaKo-S — Seed / Happy-Path Integration Run${RESET}`,
  );
  console.log(
    `${BOLD}${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${RESET}\n`,
  );

  // Shared mutable state passed through all phases
  const state: SeedState = {
    adminToken: "",
    adminUserID: "",
    teamID: "",
    translatorToken: "",
    translatorUserID: "",
    proofreaderToken: "",
    proofreaderUserID: "",
    worksetID: "",
    comicID: "",
    chapterID: "",
    pageIDs: [],
  };

  const phases: Array<{ name: string; fn: (s: SeedState) => Promise<void> }> = [
    { name: "Phase 0 — Bootstrap", fn: phase0Bootstrap },
    { name: "Phase 1 — Create Users", fn: phase1Users },
    { name: "Phase 2 — Create Content", fn: phase2Content },
    { name: "Phase 3 — Assignments", fn: phase3Assignments },
    { name: "Phase 4 — Workflow Start", fn: phase4WorkflowStart },
    { name: "Phase 5 — Translate", fn: phase5Translate },
    { name: "Phase 6 — Workflow Proofread", fn: phase6WorkflowProofread },
    { name: "Phase 7 — Proofread", fn: phase7Proofread },
    { name: "Phase 8 — Finalize", fn: phase8Finalize },
    { name: "Phase 9 — Team Board", fn: phase9TeamBoard },
  ];

  for (const { name, fn } of phases) {
    try {
      await fn(state);
    } catch (err) {
      console.error(`\n${BOLD}${RED}✗ FAILED at ${name}${RESET}`);
      console.error(`${RED}${String(err)}${RESET}\n`);
      process.exit(1);
    }
  }

  console.log(
    `\n${BOLD}${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${RESET}`,
  );
  console.log(`${BOLD}${GREEN}  ✓ All phases passed — seed complete!${RESET}`);
  console.log(
    `${BOLD}${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${RESET}\n`,
  );
}

run();

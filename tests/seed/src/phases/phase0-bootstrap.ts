/**
 * Phase 0 — Bootstrap
 *
 * 1. Login as super admin
 * 2. Create a new team
 * 3. Inject admin as member with REVIEWER + ADMIN roles
 */

import { api, logStep, logOk } from "../client";
import { SUPER_ADMIN_CREDS, ROLE } from "../config";
import type { LoginUserRes, CreateTeamRes, SeedState } from "../types";

export async function phase0Bootstrap(state: SeedState): Promise<void> {
  logStep("Phase 0", "Bootstrap — admin login, team creation, member inject");

  // ── 0.1 Login as super admin ──────────────────────────────────────────────
  const loginRes = await api<LoginUserRes>("POST", "/auth/login", {
    qq: SUPER_ADMIN_CREDS.qq,
    password: SUPER_ADMIN_CREDS.password,
  });

  state.adminToken = loginRes.access_token;
  state.adminUserID = loginRes.user_id;
  logOk("Admin logged in", { user_id: loginRes.user_id });

  // ── 0.2 Create team ───────────────────────────────────────────────────────
  const teamRes = await api<CreateTeamRes>(
    "POST",
    "/teams",
    { name: `Seed Team ${Date.now()}`, description: "Integration test team" },
    { token: state.adminToken },
  );

  state.teamID = teamRes.id;
  logOk("Team created", { team_id: teamRes.id });

  // ── 0.3 Inject admin as member (REVIEWER | ADMIN) ─────────────────────────
  await api<{ id: string }>(
    "POST",
    "/members",
    {
      team_id: state.teamID,
      user_id: state.adminUserID,
      roles:
        ROLE.RAW_PROVIDER |
        ROLE.REVIEWER |
        ROLE.ADMIN |
        ROLE.PROOFREADER |
        ROLE.TYPESETTER,
    },
    { token: state.adminToken },
  );

  logOk("Admin injected as team member (REVIEWER + ADMIN)");
}

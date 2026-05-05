/**
 * Phase 1 — Create business users
 *
 * 1. Create invitation for translator (QQ: 10002)
 * 2. Create invitation for proofreader (QQ: 10003)
 * 3. Register translator (joins team via invitation code)
 * 4. Register proofreader (joins team via invitation code)
 */

import { api, logStep, logOk, logInfo } from "../client";
import { ROLE, TRANSLATOR_CREDS, PROOFREADER_CREDS } from "../config";
import type { InvitationInfo, RegUserRes, SeedState } from "../types";

export async function phase1Users(state: SeedState): Promise<void> {
  logStep("Phase 1", "Create users — invitations + registration");

  // ── 1.1 Invitation for translator ─────────────────────────────────────────
  const invTranslator = await api<InvitationInfo>(
    "POST",
    "/member-invitations",
    {
      team_id: state.teamID,
      invitee_qid: TRANSLATOR_CREDS.qid,
      role_mask: ROLE.TRANSLATOR,
    },
    { token: state.adminToken },
  );

  logOk("Invitation created for translator", {
    code: invTranslator.invitation_code,
  });

  // ── 1.2 Invitation for proofreader ────────────────────────────────────────
  const invProofreader = await api<InvitationInfo>(
    "POST",
    "/member-invitations",
    {
      team_id: state.teamID,
      invitee_qid: PROOFREADER_CREDS.qid,
      role_mask: ROLE.PROOFREADER,
    },
    { token: state.adminToken },
  );

  logOk("Invitation created for proofreader", {
    code: invProofreader.invitation_code,
  });

  // ── 1.3 Register translator ───────────────────────────────────────────────
  const translatorRes = await api<RegUserRes>("POST", "/auth/register", {
    qid: TRANSLATOR_CREDS.qid,
    password: TRANSLATOR_CREDS.password,
    name: TRANSLATOR_CREDS.name,
    invitation_code: invTranslator.invitation_code,
  });

  state.translatorToken = translatorRes.token;
  state.translatorUserID = translatorRes.user_id;
  logOk("Translator registered", { user_id: translatorRes.user_id });

  // ── 1.4 Register proofreader ──────────────────────────────────────────────
  const proofreaderRes = await api<RegUserRes>("POST", "/auth/register", {
    qid: PROOFREADER_CREDS.qid,
    password: PROOFREADER_CREDS.password,
    name: PROOFREADER_CREDS.name,
    invitation_code: invProofreader.invitation_code,
  });

  state.proofreaderToken = proofreaderRes.token;
  state.proofreaderUserID = proofreaderRes.user_id;
  logOk("Proofreader registered", { user_id: proofreaderRes.user_id });

  logInfo("Note: Both users joined team via invitation code");
}

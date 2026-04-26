CREATE TABLE IF NOT EXISTS "t_member_invitation" (
    "id" TEXT PRIMARY KEY,
    
    "inviter_id" TEXT NOT NULL REFERENCES "t_user" ("id") ON DELETE CASCADE,
    "team_id" TEXT NOT NULL REFERENCES "t_team" ("id") ON DELETE CASCADE,

    "invitee_qid" TEXT NOT NULL,
    -- NOTE: invitation code does not have to be unique, as
    -- we only find an invitation by invitee qid.
    "invitation_code" TEXT NOT NULL,
    
    "pending" BOOLEAN NOT NULL DEFAULT TRUE,

    "role_mask" INTEGER NOT NULL DEFAULT 0,
    
    "created_at" TIMESTAMPTZ DEFAULT NOW(),
    "updated_at" TIMESTAMPTZ DEFAULT NOW(),
);

CREATE INDEX IF NOT EXISTS "idx_member_invitation_invitee_qid" ON "t_member_invitation" ("invitee_qid");
CREATE INDEX IF NOT EXISTS "idx_member_invitation_team_id_created_at_desc" ON "t_member_invitation" ("team_id", "created_at" DESC);

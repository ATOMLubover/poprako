CREATE TABLE "invitation_table" (
    "id"                   TEXT        PRIMARY KEY,

    "invitor_id"           TEXT        NOT NULL REFERENCES "user_table" ("id") ON DELETE CASCADE,
    "team_id"              TEXT        NOT NULL REFERENCES "team_table" ("id") ON DELETE CASCADE,
    "invitee_qq"           TEXT        NOT NULL,

    "invitation_code"      TEXT        UNIQUE NOT NULL,

    "to_be_raw_provider"   BOOLEAN     NOT NULL DEFAULT FALSE,
    "to_be_translator"     BOOLEAN     NOT NULL DEFAULT FALSE,
    "to_be_proofreader"    BOOLEAN     NOT NULL DEFAULT FALSE,
    "to_be_typesetter"     BOOLEAN     NOT NULL DEFAULT FALSE,
    "to_be_reviewer"       BOOLEAN     NOT NULL DEFAULT FALSE,
    "to_be_publisher"      BOOLEAN     NOT NULL DEFAULT FALSE,
    "to_be_admin"          BOOLEAN     NOT NULL DEFAULT FALSE,

    "pending"              BOOLEAN     NOT NULL DEFAULT TRUE,

    "created_at"           TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updated_at"           TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- GetByInviteeQQ(): primary lookup path
CREATE INDEX "idx_invitation_invitee_qq"
    ON "invitation_table" ("invitee_qq", "created_at" DESC);

-- List(team_id, pending): team admin views pending/resolved invitations
CREATE INDEX "idx_invitation_team_pending"
    ON "invitation_table" ("team_id", "pending", "created_at" DESC);

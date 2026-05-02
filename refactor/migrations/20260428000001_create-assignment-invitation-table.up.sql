CREATE TABLE IF NOT EXISTS "t_assignment_invitation" (
    "id" TEXT PRIMARY KEY,

    "chapter_id" TEXT NOT NULL REFERENCES "t_chapter" ("id") ON DELETE CASCADE,
    "inviter_id" TEXT NOT NULL REFERENCES "t_user" ("id") ON DELETE CASCADE,

    "invitee_qid" TEXT NOT NULL,
    "invitation_code" TEXT NOT NULL,

    "pending" BOOLEAN NOT NULL DEFAULT TRUE,
    "role_mask" BIGINT NOT NULL,

    "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS "idx_assignment_invitation_chapter_id"
    ON "t_assignment_invitation" ("chapter_id");

CREATE INDEX IF NOT EXISTS "idx_assignment_invitation_invitee_qid_pending"
    ON "t_assignment_invitation" ("invitee_qid", "pending", "created_at" DESC);

CREATE UNIQUE INDEX IF NOT EXISTS "uidx_assignment_invitation_pending_code"
    ON "t_assignment_invitation" ("invitation_code")
    WHERE "pending" = TRUE;

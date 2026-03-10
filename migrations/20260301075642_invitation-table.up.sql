CREATE TABLE "invitation_table" (
    "id"                    TEXT PRIMARY KEY,

    "invitor_id"            TEXT NOT NULL REFERENCES "user_table" ("id") ON DELETE CASCADE,
    "target_team_id"        TEXT NOT NULL REFERENCES "team_table" ("id") ON DELETE CASCADE,
    "invitee_qq"            TEXT NOT NULL,

    "invitation_code"       TEXT UNIQUE NOT NULL,

    "to_be_raw_provider"    BOOLEAN NOT NULL DEFAULT FALSE,
    "to_be_translator"      BOOLEAN NOT NULL DEFAULT FALSE,
    "to_be_proofreader"     BOOLEAN NOT NULL DEFAULT FALSE,
    "to_be_typesetter"      BOOLEAN NOT NULL DEFAULT FALSE,
    "to_be_reviewer"        BOOLEAN NOT NULL DEFAULT FALSE,
    "to_be_publisher"       BOOLEAN NOT NULL DEFAULT FALSE,
    "to_be_admin"           BOOLEAN NOT NULL DEFAULT FALSE,

    "pending"               BOOLEAN NOT NULL DEFAULT TRUE,

    "created_at"            TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updated_at"            TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX "idx_invitation_target_team_created_at_desc"
    ON "invitation_table" ("target_team_id", "created_at" DESC)
    WHERE "pending" IS FALSE;

CREATE TABLE "member_table" (
    "id"                    TEXT        PRIMARY KEY,

    "user_id"               TEXT        NOT NULL REFERENCES "user_table" ("id") ON DELETE CASCADE,
    "team_id"               TEXT        NOT NULL REFERENCES "team_table" ("id") ON DELETE CASCADE,

    "assigned_raw_provider_at" TIMESTAMPTZ,
    "assigned_translator_at"   TIMESTAMPTZ,
    "assigned_proofreader_at"  TIMESTAMPTZ,
    "assigned_typesetter_at"   TIMESTAMPTZ,
    "assigned_reviewer_at"     TIMESTAMPTZ,
    "assigned_publisher_at"    TIMESTAMPTZ,
    "assigned_admin_at"        TIMESTAMPTZ,

    "created_at"             TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updated_at"             TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "deleted_at"             TIMESTAMPTZ
);

CREATE INDEX "idx_member_user_id"
    ON "member_table" ("user_id")
    WHERE "deleted_at" IS NULL;

CREATE INDEX "idx_member_team_id"
    ON "member_table" ("team_id")
    WHERE "deleted_at" IS NULL;

CREATE INDEX "idx_member_raw_provider"
    ON "member_table" ("assigned_raw_provider_at")
    WHERE "assigned_raw_provider_at" IS NOT NULL
      AND "deleted_at" IS NULL;

CREATE INDEX "idx_member_translator"
    ON "member_table" ("assigned_translator_at")
    WHERE "assigned_translator_at" IS NOT NULL
      AND "deleted_at" IS NULL;

CREATE INDEX "idx_member_proofreader"
    ON "member_table" ("assigned_proofreader_at")
    WHERE "assigned_proofreader_at" IS NOT NULL
      AND "deleted_at" IS NULL;

CREATE INDEX "idx_member_typesetter"
    ON "member_table" ("assigned_typesetter_at")
    WHERE "assigned_typesetter_at" IS NOT NULL
      AND "deleted_at" IS NULL;

CREATE INDEX "idx_member_reviewer"
    ON "member_table" ("assigned_reviewer_at")
    WHERE "assigned_reviewer_at" IS NOT NULL
      AND "deleted_at" IS NULL;

CREATE INDEX "idx_member_publisher"
    ON "member_table" ("assigned_publisher_at")
    WHERE "assigned_publisher_at" IS NOT NULL
      AND "deleted_at" IS NULL;

CREATE INDEX "idx_member_admin"
    ON "member_table" ("assigned_admin_at")
    WHERE "assigned_admin_at" IS NOT NULL
      AND "deleted_at" IS NULL;

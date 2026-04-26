CREATE TABLE "member_table" (
    "id"                    TEXT        PRIMARY KEY,

    "user_id"               TEXT        NOT NULL REFERENCES "user_table" ("id") ON DELETE CASCADE,
    "team_id"               TEXT        NOT NULL REFERENCES "team_table" ("id") ON DELETE CASCADE,

    "assigned_raw_provider_at" TIMESTAMPTZ,
    "assigned_translator_at"   TIMESTAMPTZ,
    "assigned_proofreader_at"  TIMESTAMPTZ,
    "assigned_typesetter_at"   TIMESTAMPTZ,
    "assigned_redrawer_at"     TIMESTAMPTZ,
    "assigned_reviewer_at"     TIMESTAMPTZ,
    "assigned_publisher_at"    TIMESTAMPTZ,
    "assigned_admin_at"        TIMESTAMPTZ,

    "created_at"             TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updated_at"             TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "deleted_at"             TIMESTAMPTZ
);

-- Get/Exist(user_id, team_id): point lookup + prevents duplicate membership
CREATE UNIQUE INDEX "uidx_member_user_team"
    ON "member_table" ("user_id", "team_id")
    WHERE "deleted_at" IS NULL;

-- List(team_id): first leg of team-scoped user name search
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

-- 将预置超级管理员加入默认汉化组，并授予全部成员权限
INSERT INTO "member_table" (
    "id",
    "user_id",
    "team_id",
    "assigned_raw_provider_at",
    "assigned_translator_at",
    "assigned_proofreader_at",
    "assigned_typesetter_at",
    "assigned_redrawer_at",
    "assigned_reviewer_at",
    "assigned_publisher_at",
    "assigned_admin_at"
) VALUES (
    '00000000-0000-0000-0000-000000000201',
    '00000000-0000-0000-0000-000000000001',
    '00000000-0000-0000-0000-000000000101',
    NOW(),
    NOW(),
    NOW(),
    NOW(),
    NOW(),
    NOW(),
    NOW(),
    NOW()
);

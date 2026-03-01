CREATE TABLE "user_table" (
    "id" TEXT PRIMARY KEY,

    "name"       TEXT NOT NULL,
    "qq"         TEXT UNIQUE NOT NULL,
    "avatar_url" TEXT NOT NULL,

    "password_hash" TEXT NOT NULL,

    "assigned_picture_source_at" TIMESTAMPTZ,
    "assigned_translator_at"     TIMESTAMPTZ,
    "assigned_proofreader_at"    TIMESTAMPTZ,
    "assigned_typesetter_at"     TIMESTAMPTZ,
    "assigned_reviewer_at"       TIMESTAMPTZ,
    "assigned_admin_at"          TIMESTAMPTZ,
    "assigned_super_admin_at"    TIMESTAMPTZ,

    "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    "deleted_at" TIMESTAMPTZ
);

CREATE INDEX "idx_user_name_trgm"
    ON "user_table"
    USING GIN ("name" gin_trgm_ops)
    WHERE "deleted_at" IS NULL;

CREATE INDEX "idx_user_qq"
    ON "user_table" ("qq")
    WHERE "deleted_at" IS NULL;

CREATE INDEX "idx_user_created_at_desc"
    ON "user_table" ("created_at" DESC)
    WHERE "deleted_at" IS NULL;
CREATE INDEX "idx_user_updated_at_desc"
    ON "user_table" ("updated_at" DESC)
    WHERE "deleted_at" IS NULL;

CREATE INDEX "idx_user_picture_source"
    ON "user_table" ("assigned_picture_source_at")
    WHERE "assigned_picture_source_at" IS NOT NULL
      AND "deleted_at" IS NULL;
CREATE INDEX "idx_user_translator"
    ON "user_table" ("assigned_translator_at")
    WHERE "assigned_translator_at" IS NOT NULL
      AND "deleted_at" IS NULL;
CREATE INDEX "idx_user_proofreader"
    ON "user_table" ("assigned_proofreader_at")
    WHERE "assigned_proofreader_at" IS NOT NULL
      AND "deleted_at" IS NULL;
CREATE INDEX "idx_user_typesetter"
    ON "user_table" ("assigned_typesetter_at")
    WHERE "assigned_typesetter_at" IS NOT NULL
      AND "deleted_at" IS NULL;
CREATE INDEX "idx_user_reviewer"
    ON "user_table" ("assigned_reviewer_at")
    WHERE "assigned_reviewer_at" IS NOT NULL
      AND "deleted_at" IS NULL;
CREATE INDEX "idx_user_admin"
    ON "user_table" ("assigned_admin_at")
    WHERE "assigned_admin_at" IS NOT NULL
      AND "deleted_at" IS NULL;
CREATE INDEX "idx_user_super_admin"
    ON "user_table" ("assigned_super_admin_at")
    WHERE "assigned_super_admin_at" IS NOT NULL
      AND "deleted_at" IS NULL;

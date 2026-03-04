CREATE TABLE "user_table" (
    "id"                    TEXT PRIMARY KEY,

    "name"                  TEXT NOT NULL,
    "qq"                    TEXT UNIQUE NOT NULL,
    "avatar_url"            TEXT NOT NULL,

    "password_hash"         TEXT NOT NULL,
    
    "is_super_admin"        BOOLEAN NOT NULL DEFAULT FALSE,

    "created_at"            TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updated_at"            TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    "deleted_at"            TIMESTAMPTZ
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

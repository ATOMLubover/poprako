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

-- 预插入一个超级管理员账号，密码为 123456
INSERT INTO "user_table" (
    "id", "name", "qq", "avatar_url", "password_hash", "is_super_admin"
) VALUES (
    '00000000-0000-0000-0000-000000000001',
    '超级管理员',
    '123456789',
    '',
    '$2a$10$N9qo8uLOickgx2ZMRZo5i.ej3c3w9uQZrL6Bq4t7Yq4UstP8yW', -- 123456
    TRUE
);

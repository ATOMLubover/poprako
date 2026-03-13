CREATE TABLE "user_table" (
    "id"                 TEXT        PRIMARY KEY,
    "name"               TEXT        NOT NULL,
    "qq"                 TEXT        UNIQUE NOT NULL,

    "avatar_oss_key"     TEXT        NOT NULL,
    "is_avatar_uploaded" BOOLEAN     NOT NULL DEFAULT FALSE,

    "password_hash"      TEXT        NOT NULL,

    "is_super_admin"     BOOLEAN     NOT NULL DEFAULT FALSE,

    "created_at"         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updated_at"         TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    "deleted_at"         TIMESTAMPTZ
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
    "id",
    "name",
    "qq",
    "avatar_oss_key",
    "password_hash",
    "is_super_admin"
) VALUES (
    '00000000-0000-0000-0000-000000000001',
    'SuperAdmin OvO',
    '123456789',
    '',
    '$2a$10$eEEkAsc7h3jdkOyjahdH6OX20w/dHKdGVaH7MNREkh54O57v.E2y2', -- 123456
    TRUE
);

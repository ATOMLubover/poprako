CREATE TABLE "team_table" (
    "id"                 TEXT        PRIMARY KEY,
    "name"               TEXT        NOT NULL,
    "description"        TEXT,

    "avatar_oss_key"     TEXT,
    "is_avatar_uploaded" BOOLEAN     NOT NULL DEFAULT FALSE,

    "created_at"         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updated_at"         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "deleted_at"         TIMESTAMPTZ
);

-- Unique team name constraint + List() lookup
CREATE UNIQUE INDEX "uidx_team_name"
    ON "team_table" ("name")
    WHERE "deleted_at" IS NULL;

-- 预插入一个默认汉化组，供超级管理员加入并管理系统内初始资源
INSERT INTO "team_table" (
    "id",
    "name",
    "description",
    "avatar_oss_key",
    "is_avatar_uploaded"
) VALUES (
    '00000000-0000-0000-0000-000000000101',
    '默认汉化组',
    '系统初始化时自动创建的默认汉化组。',
    '',
    FALSE
);

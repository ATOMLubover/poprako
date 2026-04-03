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
CREATE TABLE "comic_table" (
    "id"                    TEXT PRIMARY KEY,

    "team_id"               TEXT NOT NULL REFERENCES "team_table" ("id") ON DELETE CASCADE,

    "index"                 INTEGER NOT NULL DEFAULT 0,
    "title"                 TEXT NOT NULL,
    "author"                TEXT NOT NULL,
    "description"           TEXT,

    "cover_url"             TEXT,
    "chapter_count"         INTEGER NOT NULL DEFAULT 0,

    "creator_id"            TEXT NOT NULL REFERENCES "user_table" ("id") ON DELETE RESTRICT,

    "created_at"            TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updated_at"            TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "deleted_at"            TIMESTAMPTZ
);

CREATE INDEX "idx_comic_team_created_at_desc"
    ON "comic_table" ("team_id", "created_at" DESC)
    WHERE "deleted_at" IS NULL;

CREATE INDEX "idx_comic_creator_id"
    ON "comic_table" ("creator_id")
    WHERE "deleted_at" IS NULL;

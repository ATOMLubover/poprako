CREATE TABLE IF NOT EXISTS "t_comic" (
    "id"                TEXT        PRIMARY KEY,

    "workset_id"        TEXT        NOT NULL REFERENCES "t_workset" ("id") ON DELETE CASCADE,

    "index"             INTEGER     NOT NULL,
    "title"             TEXT        NOT NULL,
    "author"            TEXT        NOT NULL,
    "fuzzy_title"       TEXT        NOT NULL,
    "description"       TEXT,

    "is_completed"      BOOLEAN     NOT NULL DEFAULT FALSE,

    "cover_key"         TEXT,
    "cover_uploaded"    BOOLEAN     NOT NULL DEFAULT FALSE,

    -- chapter_count is the count of **active** chapters.
    "chapter_count"     INTEGER     NOT NULL DEFAULT 0,
    -- chapter_next_index is the next index to be assigned to a new chapter, 
    -- It is calculated every time a new chapter is added, and is not affected by chapter deletions.
    "chapter_next_index"  INTEGER     NOT NULL DEFAULT 0,

    "creator_id"        TEXT        NOT NULL REFERENCES "t_user" ("id"),

    "last_active_at"    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "created_at"        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updated_at"        TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Unique index ensures comics use unique index inside one workset.
CREATE UNIQUE INDEX IF NOT EXISTS "uidx_comic_workset_id_index"
    ON "t_comic" ("workset_id", "index");

-- Support fuzzy search on dedicated fuzzy title
CREATE INDEX IF NOT EXISTS "trgm_idx_comic_fuzzy_title"
    ON "t_comic" USING gin ("fuzzy_title" gin_trgm_ops)
    WHERE is_completed = FALSE;

-- Support title-only fuzzy search with last active sort under one workset
CREATE INDEX IF NOT EXISTS "idx_comic_workset_id_last_active"
    ON "t_comic" ("workset_id", "last_active_at" DESC);

-- Support filtered comic search that excludes completed rows.
CREATE INDEX IF NOT EXISTS "idx_comic_workset_id_last_active_filtered"
    ON "t_comic" ("workset_id", "last_active_at" DESC)
    WHERE is_completed = FALSE;

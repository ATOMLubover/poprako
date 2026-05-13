CREATE TABLE IF NOT EXISTS "t_chapter" (
    "id"                    TEXT        PRIMARY KEY,

    "comic_id"              TEXT        NOT NULL REFERENCES "t_comic" ("id") ON DELETE CASCADE,
    "pinned"                BOOLEAN     NOT NULL DEFAULT FALSE,

    "index"                 INTEGER     NOT NULL,
    "subtitle"              TEXT        NOT NULL,

    "page_count"            INTEGER     NOT NULL DEFAULT 0,
    "total_unit_count"      INTEGER     NOT NULL DEFAULT 0,
    "translated_unit_count" INTEGER     NOT NULL DEFAULT 0,
    "proofread_unit_count"  INTEGER     NOT NULL DEFAULT 0,

    "uploaded_at"           TIMESTAMPTZ,
    "translating_at"       TIMESTAMPTZ,
    "translated_at"         TIMESTAMPTZ,
    "proofreading_at"       TIMESTAMPTZ,
    "proofread_at"          TIMESTAMPTZ,
    "typesetting_at"        TIMESTAMPTZ,
    "typeset_at"            TIMESTAMPTZ,
    "reviewed_at"           TIMESTAMPTZ,
    "published_at"          TIMESTAMPTZ,

    "creator_id"            TEXT        NOT NULL REFERENCES "t_user" ("id"),

    "created_at"            TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updated_at"            TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Unique position per comic.
CREATE UNIQUE INDEX IF NOT EXISTS "uidx_chapter_comic_id_index"
    ON "t_chapter" ("comic_id", "index");

-- Support chapter list query by comic and index desc
CREATE INDEX IF NOT EXISTS "idx_chapter_comic_id_index_desc"
    ON "t_chapter" ("comic_id", "index" DESC);

-- Enforce single pinned chapter per comic while allowing many unpinned rows.
CREATE UNIQUE INDEX IF NOT EXISTS "uidx_chapter_comic_id_pinned_true"
    ON "t_chapter" ("comic_id")
    WHERE pinned = TRUE;

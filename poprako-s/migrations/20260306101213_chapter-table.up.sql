CREATE TABLE "chapter_table" (
    "id"                    TEXT        PRIMARY KEY,

    "comic_id"              TEXT        NOT NULL REFERENCES "comic_table" ("id") ON DELETE CASCADE,
    "pinned"             BOOLEAN     NOT NULL DEFAULT FALSE,

    "index"                 INTEGER     NOT NULL DEFAULT 0,
    "subtitle"              TEXT        NOT NULL,

    "page_count"            INTEGER     NOT NULL DEFAULT 0,
    "total_unit_count"      INTEGER     NOT NULL DEFAULT 0,
    "translated_unit_count" INTEGER     NOT NULL DEFAULT 0,
    "proofread_unit_count"  INTEGER     NOT NULL DEFAULT 0,

    "uploaded_at"           TIMESTAMPTZ,
    "transalating_at"       TIMESTAMPTZ,
    "translated_at"         TIMESTAMPTZ,
    "proofreading_at"       TIMESTAMPTZ,
    "proofread_at"          TIMESTAMPTZ,
    "typesetting_at"        TIMESTAMPTZ,
    "typeset_at"            TIMESTAMPTZ,
    "reviewed_at"           TIMESTAMPTZ,
    "published_at"          TIMESTAMPTZ,

    "creator_id"            TEXT        NOT NULL REFERENCES "user_table" ("id") ON DELETE RESTRICT,

    "created_at"            TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updated_at"            TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "deleted_at"            TIMESTAMPTZ
);

CREATE UNIQUE INDEX "uidx_chapter_comic_id_index"
    ON "chapter_table" ("comic_id", "index" DESC)
    WHERE "deleted_at" IS NULL;

CREATE INDEX "idx_chapter_comic_id"
    ON "chapter_table" ("comic_id")
    WHERE "deleted_at" IS NULL;

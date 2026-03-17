CREATE TABLE "comic_table" (
    "id"                 TEXT        PRIMARY KEY,

    "workset_id"         TEXT        NOT NULL REFERENCES "workset_table" ("id") ON DELETE CASCADE,

    "index"              INTEGER     NOT NULL DEFAULT 0,
    "title"              TEXT        NOT NULL,
    "author"             TEXT        NOT NULL,
    "description"        TEXT,

    "chapter_count"      INTEGER     NOT NULL DEFAULT 0,

    "creator_id"         TEXT        NOT NULL REFERENCES "user_table" ("id") ON DELETE RESTRICT,

    "last_active_at"     TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    "created_at"         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updated_at"         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "deleted_at"         TIMESTAMPTZ
);

CREATE UNIQUE INDEX "uidx_comic_workset_id_index"
    ON "comic_table" ("workset_id", "index")
    WHERE "deleted_at" IS NULL;

CREATE INDEX "idx_comic_workset_id_created_at_desc"
    ON "comic_table" ("workset_id", "created_at" DESC)
    WHERE "deleted_at" IS NULL;

CREATE INDEX "idx_comic_creator_id"
    ON "comic_table" ("creator_id")
    WHERE "deleted_at" IS NULL;

CREATE INDEX "idx_comic_last_active_at_desc"
    ON "comic_table" ("workset_id", "last_active_at" DESC)
    WHERE "deleted_at" IS NULL;

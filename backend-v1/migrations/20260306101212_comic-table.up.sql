CREATE TABLE "comic_table" (
    "id"                 TEXT        PRIMARY KEY,

    "workset_id"         TEXT        NOT NULL REFERENCES "workset_table" ("id") ON DELETE CASCADE,

    "index"              INTEGER     NOT NULL DEFAULT 0,
    "title"              TEXT        NOT NULL,
    "author"             TEXT        NOT NULL,
    "composed_title"     TEXT        NOT NULL,
    "description"        TEXT,

    "chapter_count"      INTEGER     NOT NULL DEFAULT 0,

    "latest_uploaded_at"     TIMESTAMPTZ,
    "latest_transalating_at" TIMESTAMPTZ,
    "latest_translated_at"   TIMESTAMPTZ,
    "latest_proofreading_at" TIMESTAMPTZ,
    "latest_proofread_at"    TIMESTAMPTZ,
    "latest_typesetting_at"  TIMESTAMPTZ,
    "latest_typeset_at"      TIMESTAMPTZ,
    "latest_reviewed_at"     TIMESTAMPTZ,
    "latest_published_at"    TIMESTAMPTZ,

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

CREATE INDEX "idx_comic_latest_uploaded_at"
    ON "comic_table" ("workset_id", "latest_uploaded_at")
    WHERE "deleted_at" IS NULL;

CREATE INDEX "idx_comic_latest_transalating_translated_at"
    ON "comic_table" ("workset_id", "latest_transalating_at", "latest_translated_at")
    WHERE "deleted_at" IS NULL;

CREATE INDEX "idx_comic_latest_proofreading_proofread_at"
    ON "comic_table" ("workset_id", "latest_proofreading_at", "latest_proofread_at")
    WHERE "deleted_at" IS NULL;

CREATE INDEX "idx_comic_latest_typesetting_typeset_at"
    ON "comic_table" ("workset_id", "latest_typesetting_at", "latest_typeset_at")
    WHERE "deleted_at" IS NULL;

CREATE INDEX "idx_comic_latest_reviewed_at"
    ON "comic_table" ("workset_id", "latest_reviewed_at")
    WHERE "deleted_at" IS NULL;

CREATE INDEX "idx_comic_latest_published_at"
    ON "comic_table" ("workset_id", "latest_published_at")
    WHERE "deleted_at" IS NULL;

CREATE INDEX "idx_comic_composed_title_trgm"
    ON "comic_table"
    USING gin (composed_title gin_trgm_ops)
    WHERE "deleted_at" IS NULL;

CREATE TABLE "comic_table" (
    "id"                 TEXT        PRIMARY KEY,

    "workset_id"         TEXT        NOT NULL REFERENCES "workset_table" ("id") ON DELETE CASCADE,

    "index"              INTEGER     NOT NULL DEFAULT 0,
    "title"              TEXT        NOT NULL,
    "author"             TEXT        NOT NULL,
    "composed_title"     TEXT        NOT NULL,
    "description"        TEXT,

    "chapter_count"      INTEGER     NOT NULL DEFAULT 0,
    
    "has_pinned_chapter"         BOOLEAN     NOT NULL DEFAULT FALSE,
    "pinned_uploaded_at"         TIMESTAMPTZ,
    "pinned_transalating_at"     TIMESTAMPTZ,
    "pinned_translated_at"       TIMESTAMPTZ,
    "pinned_proofreading_at"     TIMESTAMPTZ,
    "pinned_proofread_at"        TIMESTAMPTZ,
    "pinned_typesetting_at"      TIMESTAMPTZ,
    "pinned_typeset_at"          TIMESTAMPTZ,
    "pinned_reviewed_at"         TIMESTAMPTZ,
    "pinned_published_at"        TIMESTAMPTZ,

    "creator_id"         TEXT        NOT NULL REFERENCES "user_table" ("id") ON DELETE RESTRICT,
    "last_active_at"     TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    "created_at"         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updated_at"         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "deleted_at"         TIMESTAMPTZ
);

-- Unique position per workset (soft-delete aware)
CREATE UNIQUE INDEX "uidx_comic_workset_id_index"
    ON "comic_table" ("workset_id", "index")
    WHERE "deleted_at" IS NULL;

-- List() base scan: workset_id filter + ORDER BY index ASC
CREATE INDEX "idx_comic_workset_id_index"
    ON "comic_table" ("workset_id", "index" ASC)
    WHERE "deleted_at" IS NULL;

-- List() ORDER BY last_active_at DESC
CREATE INDEX "idx_comic_workset_last_active"
    ON "comic_table" ("workset_id", "last_active_at" DESC)
    WHERE "deleted_at" IS NULL;

-- List(fuzzy_title): GIN trigram for ILIKE search on composed_title
CREATE INDEX "idx_comic_composed_title_trgm"
    ON "comic_table"
    USING GIN ("composed_title" gin_trgm_ops)
    WHERE "deleted_at" IS NULL;

-- Upload: has_pinned_chapter=FALSE OR pinned_uploaded_at IS NULL
CREATE INDEX "idx_comic_pinned_upload"
    ON "comic_table" ("workset_id", "has_pinned_chapter", "pinned_uploaded_at")
    WHERE "deleted_at" IS NULL;

-- Translate: pinned_transalating_at IS NOT NULL [AND pinned_translated_at IS NULL]
CREATE INDEX "idx_comic_pinned_translate"
    ON "comic_table" ("workset_id", "pinned_transalating_at", "pinned_translated_at")
    WHERE "deleted_at" IS NULL;

-- Proofread: pinned_proofreading_at IS NOT NULL [AND pinned_proofread_at IS NULL]
CREATE INDEX "idx_comic_pinned_proofread"
    ON "comic_table" ("workset_id", "pinned_proofreading_at", "pinned_proofread_at")
    WHERE "deleted_at" IS NULL;

-- Typeset: pinned_typesetting_at IS NOT NULL [AND pinned_typeset_at IS NULL]
CREATE INDEX "idx_comic_pinned_typeset"
    ON "comic_table" ("workset_id", "pinned_typesetting_at", "pinned_typeset_at")
    WHERE "deleted_at" IS NULL;

-- Review: pinned_reviewed_at IS NOT NULL
CREATE INDEX "idx_comic_pinned_review"
    ON "comic_table" ("workset_id", "pinned_reviewed_at")
    WHERE "deleted_at" IS NULL;

-- Publish: pinned_published_at IS NOT NULL
CREATE INDEX "idx_comic_pinned_publish"
    ON "comic_table" ("workset_id", "pinned_published_at")
    WHERE "deleted_at" IS NULL;

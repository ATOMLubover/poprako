CREATE TABLE "assignment_table" (
    "id"                    TEXT PRIMARY KEY,

    "comic_id"              TEXT NOT NULL REFERENCES "comic_table" ("id") ON DELETE CASCADE,
    "user_id"               TEXT NOT NULL REFERENCES "user_table" ("id") ON DELETE CASCADE,
    
    "assigned_translator_at"    TIMESTAMPTZ,
    "assigned_proofreader_at"   TIMESTAMPTZ,
    "assigned_typesetter_at"    TIMESTAMPTZ,
    "assigned_reviewer_at"      TIMESTAMPTZ,
    "assigned_uploader_at"      TIMESTAMPTZ,

    "created_at"            TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updated_at"            TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "deleted_at"            TIMESTAMPTZ
);

CREATE INDEX "idx_assignment_comic_id"
    ON "assignment_table" ("comic_id")
    WHERE "deleted_at" IS NULL;
CREATE INDEX "idx_assignment_user_id"
    ON "assignment_table" ("user_id")
    WHERE "deleted_at" IS NULL;

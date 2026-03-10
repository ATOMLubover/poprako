CREATE TABLE "assignment_table" (
    "id"                    TEXT PRIMARY KEY,

    "chapter_id"              TEXT NOT NULL REFERENCES "chapter_table" ("id") ON DELETE CASCADE,
    "user_id"               TEXT NOT NULL REFERENCES "user_table" ("id") ON DELETE CASCADE,
    
    "assigned_raw_provider_at"  TIMESTAMPTZ,
    "assigned_translator_at"    TIMESTAMPTZ,
    "assigned_proofreader_at"   TIMESTAMPTZ,
    "assigned_typesetter_at"    TIMESTAMPTZ,
    "assigned_reviewer_at"      TIMESTAMPTZ,
    "assigned_publisher_at"     TIMESTAMPTZ,

    "created_at"            TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updated_at"            TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX "idx_assignment_chapter_id"
    ON "assignment_table" ("chapter_id");
CREATE INDEX "idx_assignment_user_id"
    ON "assignment_table" ("user_id");

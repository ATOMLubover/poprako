CREATE TABLE IF NOT EXISTS "t_assignment" (
    "id" TEXT PRIMARY KEY,

    "chapter_id" TEXT NOT NULL REFERENCES "t_chapter" ("id") ON DELETE CASCADE,
    "user_id" TEXT NOT NULL REFERENCES "t_user" ("id") ON DELETE CASCADE,

    "assigned_raw_provider_at" TIMESTAMPTZ,
    "assigned_translator_at" TIMESTAMPTZ,
    "assigned_proofreader_at" TIMESTAMPTZ,
    "assigned_typesetter_at" TIMESTAMPTZ,
    "assigned_redrawer_at" TIMESTAMPTZ,
    "assigned_reviewer_at" TIMESTAMPTZ,
    "assigned_publisher_at" TIMESTAMPTZ,

    "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE ("chapter_id", "user_id")
);

CREATE INDEX IF NOT EXISTS "idx_assignment_chapter_id"
    ON "t_assignment" ("chapter_id", "created_at" DESC);

CREATE INDEX IF NOT EXISTS "idx_assignment_user_id"
    ON "t_assignment" ("user_id", "created_at" DESC);

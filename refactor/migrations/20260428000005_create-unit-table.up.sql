CREATE TABLE IF NOT EXISTS "t_unit" (
    "id" TEXT PRIMARY KEY,

    "page_id" TEXT NOT NULL REFERENCES "t_page" ("id") ON DELETE CASCADE,
    "index" INTEGER NOT NULL,
    "is_bubble" BOOLEAN NOT NULL DEFAULT FALSE,
    "is_proofread" BOOLEAN NOT NULL DEFAULT FALSE,

    "x_coord" DOUBLE NOT NULL,
    "y_coord" DOUBLE NOT NULL,

    "translated_text" TEXT,
    "translator_comment" TEXT,
    "last_translator_id" TEXT REFERENCES "t_user" ("id") ON DELETE SET NULL,

    "proofread_text" TEXT,
    "proofreader_comment" TEXT,
    "last_proofreader_id" TEXT REFERENCES "t_user" ("id") ON DELETE SET NULL,

    "created_at" TIMESTAMPTZ DEFAULT NOW(),
    "updated_at" TIMESTAMPTZ DEFAULT NOW(),

    UNIQUE ("page_id", "index")
);

CREATE INDEX IF NOT EXISTS "idx_unit_page_id" ON "t_unit" ("page_id");

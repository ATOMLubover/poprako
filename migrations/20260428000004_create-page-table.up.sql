CREATE TABLE IF NOT EXISTS "t_page" (
    "id" TEXT PRIMARY KEY,

    "chapter_id" TEXT NOT NULL REFERENCES "t_chapter" ("id") ON DELETE CASCADE,
    "index" INTEGER NOT NULL,

    "image_key" TEXT,
    "image_uploaded" BOOLEAN NOT NULL DEFAULT FALSE,

    "total_unit_count" INTEGER NOT NULL DEFAULT 0,
    "translated_unit_count" INTEGER NOT NULL DEFAULT 0,
    "proofread_unit_count" INTEGER NOT NULL DEFAULT 0,

    "created_at" TIMESTAMPTZ DEFAULT NOW(),
    "updated_at" TIMESTAMPTZ DEFAULT NOW(),

    UNIQUE ("chapter_id", "index")
);

CREATE INDEX IF NOT EXISTS "idx_page_chapter_id" ON "t_page" ("chapter_id");


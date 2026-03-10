CREATE TABLE "page_table" (
    "id"                    TEXT PRIMARY KEY,

    "chapter_id"            TEXT NOT NULL REFERENCES "chapter_table" ("id") ON DELETE CASCADE,

    "index"                 INTEGER NOT NULL DEFAULT 0,
    "oss_key"               TEXT,
    
    "uploaded"              BOOLEAN NOT NULL DEFAULT FALSE,
    
    "total_unit_count"      INTEGER NOT NULL DEFAULT 0,
    "translated_unit_count" INTEGER NOT NULL DEFAULT 0,
    "proved_unit_count"     INTEGER NOT NULL DEFAULT 0,
    
    "creator_id"            TEXT NOT NULL REFERENCES "user_table" ("id") ON DELETE RESTRICT,

    "created_at"            TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updated_at"            TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX "idx_page_chapter_id"
    ON "page_table" ("chapter_id");

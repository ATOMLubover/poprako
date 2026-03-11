ALTER TABLE "chapter_table"
    DROP COLUMN IF EXISTS "uploaded_at",
    DROP COLUMN IF EXISTS "transalating_at",
    DROP COLUMN IF EXISTS "translated_at",
    DROP COLUMN IF EXISTS "proofreading_at",
    DROP COLUMN IF EXISTS "proofread_at",
    DROP COLUMN IF EXISTS "typesetting_at",
    DROP COLUMN IF EXISTS "typeset_at",
    DROP COLUMN IF EXISTS "reviewed_at",
    DROP COLUMN IF EXISTS "published_at";

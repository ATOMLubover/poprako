ALTER TABLE "chapter_table"
    ADD COLUMN "uploaded_at"      TIMESTAMPTZ,
    ADD COLUMN "transalating_at"  TIMESTAMPTZ,
    ADD COLUMN "translated_at"    TIMESTAMPTZ,
    ADD COLUMN "proofreading_at"  TIMESTAMPTZ,
    ADD COLUMN "proofread_at"     TIMESTAMPTZ,
    ADD COLUMN "typesetting_at"   TIMESTAMPTZ,
    ADD COLUMN "typeset_at"       TIMESTAMPTZ,
    ADD COLUMN "reviewed_at"      TIMESTAMPTZ,
    ADD COLUMN "published_at"     TIMESTAMPTZ;

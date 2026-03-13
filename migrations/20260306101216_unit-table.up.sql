CREATE TABLE "unit_table" (
    "id"                 TEXT        PRIMARY KEY,

    "page_id"            TEXT        NOT NULL REFERENCES "page_table" ("id") ON DELETE CASCADE,

    "x_coord"            REAL        NOT NULL,
    "y_coord"            REAL        NOT NULL,

    "index"              INTEGER     NOT NULL,
    "in_bubble"          BOOLEAN     NOT NULL DEFAULT TRUE,

    "is_proofread"       BOOLEAN     NOT NULL DEFAULT FALSE,

    "translated_text"    TEXT,
    "translator_id"      TEXT        REFERENCES "user_table" ("id") ON DELETE SET NULL,
    "translator_comment" TEXT,

    "proofreader_text"   TEXT,
    "proofreader_id"     TEXT        REFERENCES "user_table" ("id") ON DELETE SET NULL,
    "proofreader_comment" TEXT,

    "created_at"         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updated_at"         TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE ("page_id", "index")
);

CREATE INDEX "unit_table_page_id_index"
    ON "unit_table" ("page_id");

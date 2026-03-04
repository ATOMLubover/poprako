CREATE TABLE "team_table" (
    "id"                    TEXT PRIMARY KEY,

    "name"                  TEXT NOT NULL,
    "description"           TEXT,

    "created_at"            TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updated_at"            TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "deleted_at"            TIMESTAMPTZ
);
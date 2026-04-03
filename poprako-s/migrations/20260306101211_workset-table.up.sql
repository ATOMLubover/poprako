CREATE TABLE "workset_table" (
    "id"           TEXT        PRIMARY KEY,

    "team_id"      TEXT        NOT NULL REFERENCES "team_table" ("id") ON DELETE CASCADE,
    "index"        INTEGER     NOT NULL,

    "name"         TEXT        NOT NULL,
    "description"  TEXT,
    "comic_count"  INTEGER     NOT NULL DEFAULT 0,

    "created_at"   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updated_at"   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- List(team_id) ORDER BY index ASC; unique position per team
CREATE UNIQUE INDEX "uidx_workset_team_id_index"
    ON "workset_table" ("team_id", "index");
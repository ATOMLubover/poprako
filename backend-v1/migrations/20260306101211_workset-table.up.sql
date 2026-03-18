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

CREATE INDEX "idx_workset_team_id"
    ON "workset_table" ("team_id");

CREATE UNIQUE INDEX "uidx_workset_team_id_index"
    ON "workset_table" ("team_id", "index");
CREATE TABLE "workset_table" (
    "id"           SERIAL PRIMARY KEY,

    "team_id"      INTEGER NOT NULL REFERENCES "team_table" ("id") ON DELETE CASCADE,
    "index"        INTEGER NOT NULL,

    "name"         VARCHAR(255) NOT NULL,
    "description"  TEXT,
    "comic_count"  INTEGER NOT NULL DEFAULT 0,

    "created_at"   TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    "updated_at"   TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX "idx_workset_team_id"
    ON "workset_table" ("team_id");

CREATE UNIQUE INDEX "uidx_workset_team_id_index"
    ON "workset_table" ("team_id", "index");
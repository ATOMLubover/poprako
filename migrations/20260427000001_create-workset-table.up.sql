CREATE TABLE IF NOT EXISTS "t_workset" (
    "id"          TEXT        PRIMARY KEY,

    "team_id"     TEXT        NOT NULL REFERENCES "t_team" ("id") ON DELETE CASCADE,
    "index"       INTEGER     NOT NULL,

    "name"        TEXT        NOT NULL,
    "description" TEXT,

    -- comic_count is the count of **active** comics in the workset.
    "comic_count" INTEGER     NOT NULL DEFAULT 0,
    -- comic_next_index is the next index to assign to a new comic in the workset.
    -- It is calculated every time a new comic is added, and is not affected by comic deletions.
    "comic_next_index" INTEGER     NOT NULL DEFAULT 0,

    "created_at"  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updated_at"  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Unique index ensures per-team workset index uniqueness.
CREATE UNIQUE INDEX IF NOT EXISTS "uidx_workset_team_id_index"
    ON "t_workset" ("team_id", "index");

-- Support fast team-scoped list queries.
CREATE INDEX IF NOT EXISTS "idx_workset_team_id"
    ON "t_workset" ("team_id");

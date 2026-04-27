CREATE TABLE IF NOT EXISTS "t_workset" (
    "id"          TEXT        PRIMARY KEY,

    "team_id"     TEXT        NOT NULL REFERENCES "t_team" ("id") ON DELETE CASCADE,
    "index"       INTEGER     NOT NULL,

    "name"        TEXT        NOT NULL,
    "description" TEXT,
    "comic_count" INTEGER     NOT NULL DEFAULT 0,

    "created_at"  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updated_at"  TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- Soft-delete timestamp. NULL means the workset is active.
    "deleted_at"  TIMESTAMPTZ
);

-- Partial unique index: only active (non-deleted) worksets must have a
-- unique index within a team. Deleted worksets are excluded so that
-- their stale index values do not block future inserts.
CREATE UNIQUE INDEX IF NOT EXISTS "uidx_workset_team_id_index"
    ON "t_workset" ("team_id", "index")
    WHERE deleted_at IS NULL;

-- Support fast team-scoped list queries.
CREATE INDEX IF NOT EXISTS "idx_workset_team_id"
    ON "t_workset" ("team_id");

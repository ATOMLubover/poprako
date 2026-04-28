CREATE TABLE IF NOT EXISTS "t_user_stats" (
    "id" TEXT PRIMARY KEY,

    "user_id" TEXT NOT NULL REFERENCES "t_user" ("id") ON DELETE CASCADE,

    "total_assignment_count" INTEGER NOT NULL DEFAULT 0,
    "active_assignment_count" INTEGER NOT NULL DEFAULT 0,
    "finished_assignment_count" INTEGER NOT NULL DEFAULT 0,

    "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS "uidx_user_stats_user_id"
    ON "t_user_stats" ("user_id");

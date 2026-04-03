CREATE TABLE "user_stats_table" (
    "id" TEXT PRIMARY KEY,

    "user_id" TEXT NOT NULL REFERENCES "user_table" ("id") ON DELETE CASCADE,

    "total_assignment_count" INTEGER NOT NULL DEFAULT 0,
    "active_assignment_count" INTEGER NOT NULL DEFAULT 0,
    "finished_assignment_count" INTEGER NOT NULL DEFAULT 0,
    
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- One-to-one with user_table; UNIQUE enforces and optimises GetOrCreateStats
CREATE UNIQUE INDEX "uidx_user_stats_user_id"
    ON "user_stats_table" ("user_id");


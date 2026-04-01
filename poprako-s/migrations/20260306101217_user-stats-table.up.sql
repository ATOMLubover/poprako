CREATE TABLE "user_stats_table" (
    "id" TEXT PRIMARY KEY,

    "user_id" TEXT NOT NULL REFERENCES "user_table" ("id") ON DELETE CASCADE,

    "total_assignment_count" INTEGER NOT NULL DEFAULT 0,
    "active_assignment_count" INTEGER NOT NULL DEFAULT 0,
    "finished_assignment_count" INTEGER NOT NULL DEFAULT 0,
    
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX "idx_user_stats_table_user_id" ON "user_stats_table" ("user_id");


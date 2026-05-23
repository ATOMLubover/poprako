CREATE TABLE IF NOT EXISTS "t_announcement" (
    "id" TEXT PRIMARY KEY,

    "team_id" TEXT NOT NULL REFERENCES "t_team"("id") ON DELETE CASCADE,
    "user_id" TEXT NOT NULL REFERENCES "t_user"("id") ON DELETE CASCADE,

    "title" TEXT NOT NULL,
    "content" TEXT NOT NULL,

    "created_at" TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS "idx_announcement_team_id_created_at_desc" ON "t_announcement" ("team_id", "created_at" DESC);


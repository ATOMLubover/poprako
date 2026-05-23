CREATE TABLE IF NOT EXISTS "t_comment" (
    "id" TEXT PRIMARY KEY,

    "team_id" TEXT NOT NULL REFERENCES "t_team"("id") ON DELETE CASCADE, 
    "user_id" TEXT NOT NULL REFERENCES "t_user"("id") ON DELETE CASCADE,

    "content" TEXT NOT NULL,

    "created_at" TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS "idx_comment_team_id_created_at_desc" ON "t_comment" ("team_id", "created_at" DESC);

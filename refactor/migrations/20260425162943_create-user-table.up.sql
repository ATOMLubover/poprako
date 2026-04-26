CREATE TABLE IF NOT EXISTS "t_user" (
    "id" TEXT PRIMARY KEY,

    "qid" TEXT NOT NULL UNIQUE,
    "nickname" VARCHAR(255) NOT NULL UNIQUE,

    "password_hash" VARCHAR(255) NOT NULL,

    "last_active_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    "created_at" TIMESTAMPTZ DEFAULT NOW(),
    "updated_at" TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS "idx_user_qid" ON "t_user" ("qid");
CREATE INDEX IF NOT EXISTS "trgm_idx_user_nickname" ON "t_user" USING gin ("nickname" gin_trgm_ops);

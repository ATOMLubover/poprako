CREATE TABLE IF NOT EXISTS "t_user" (
    "id" TEXT PRIMARY KEY,

    "qid" TEXT NOT NULL UNIQUE,
    "nickname" TEXT NOT NULL UNIQUE,
    
    "avatar_key" TEXT,
    "avatar_uploaded" BOOLEAN DEFAULT FALSE,

    "password_hash" TEXT NOT NULL,

    "is_super_admin" BOOLEAN DEFAULT FALSE,
    "last_active_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    "created_at" TIMESTAMPTZ DEFAULT NOW(),
    "updated_at" TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS "idx_user_qid"
  ON "t_user" ("qid");
CREATE INDEX IF NOT EXISTS "trgm_idx_user_nickname"
  ON "t_user" USING gin ("nickname" gin_trgm_ops);
CREATE INDEX IF NOT EXISTS "idx_user_qid"
  ON "t_user" ("qid");

-- Create super admin directly in database.
INSERT INTO "t_user" (
    "id", 
    "qid", 
    "nickname", 
    "password",
    "is_super_admin"
) VALUES (
    'user-00000000-0000-0000-0000-000000000001',
    'SuperAdmin-OvO',
    '123456789',
    '$2a$10$eEEkAsc7h3jdkOyjahdH6OX20w/dHKdGVaH7MNREkh54O57v.E2y2', -- 123456
    TRUE
) ON CONFLICT (id) DO NOTHING;

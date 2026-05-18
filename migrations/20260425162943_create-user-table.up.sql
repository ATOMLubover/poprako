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

-- Create super admin directly in database.
INSERT INTO "t_user" (
    "id", 
    "qid", 
    "nickname", 
    "password_hash",
    "is_super_admin"
) VALUES (
    'user-00000000-0000-0000-0000-000000000001',
    '123456789',
    'SuperAdmin-OvO',
    '$2a$10$eEEkAsc7h3jdkOyjahdH6OX20w/dHKdGVaH7MNREkh54O57v.E2y2', -- 123456
    TRUE
) ON CONFLICT (id) DO NOTHING;

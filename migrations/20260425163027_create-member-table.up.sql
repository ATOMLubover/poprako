CREATE TABLE IF NOT EXISTS "t_member" (
    "id" TEXT PRIMARY KEY,
    
    "user_id" TEXT NOT NULL REFERENCES "t_user" ("id") ON DELETE CASCADE,
    "user_nickname" TEXT NOT NULL,
    "team_id" TEXT NOT NULL REFERENCES "t_team" ("id") ON DELETE CASCADE,
    
    "assigned_raw_provider_at" TIMESTAMPTZ,
    "assigned_translator_at" TIMESTAMPTZ,
    "assigned_proofreader_at" TIMESTAMPTZ,
    "assigned_typesetter_at" TIMESTAMPTZ,
    "assigned_redrawer_at" TIMESTAMPTZ,
    "assigned_reviewer_at" TIMESTAMPTZ,
    "assigned_publisher_at" TIMESTAMPTZ,
    "assigned_admin_at" TIMESTAMPTZ,
    
    "created_at" TIMESTAMPTZ DEFAULT NOW(),
    "updated_at" TIMESTAMPTZ DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS "uidx_member_user_team"
    ON "t_member" ("user_id", "team_id");
CREATE INDEX IF NOT EXISTS "idx_member_user_id"
    ON "t_member" ("user_id");
CREATE INDEX IF NOT EXISTS "idx_member_team_id"
    ON "t_member" ("team_id");

CREATE INDEX IF NOT EXISTS "idx_member_team_raw_provider"
    ON "t_member" ("team_id")
    WHERE "assigned_raw_provider_at" IS NOT NULL;
CREATE INDEX IF NOT EXISTS "idx_member_team_translator"
    ON "t_member" ("team_id")
    WHERE "assigned_translator_at" IS NOT NULL;
CREATE INDEX IF NOT EXISTS "idx_member_team_proofreader"
    ON "t_member" ("team_id")
    WHERE "assigned_proofreader_at" IS NOT NULL;
CREATE INDEX IF NOT EXISTS "idx_member_team_typesetter"
    ON "t_member" ("team_id")
    WHERE "assigned_typesetter_at" IS NOT NULL;
CREATE INDEX IF NOT EXISTS "idx_member_team_redrawer"
    ON "t_member" ("team_id")
    WHERE "assigned_redrawer_at" IS NOT NULL;
CREATE INDEX IF NOT EXISTS "idx_member_team_reviewer"
    ON "t_member" ("team_id")
    WHERE "assigned_reviewer_at" IS NOT NULL;
CREATE INDEX IF NOT EXISTS "idx_member_team_publisher"
    ON "t_member" ("team_id")
    WHERE "assigned_publisher_at" IS NOT NULL;
CREATE INDEX IF NOT EXISTS "idx_member_team_admin"
    ON "t_member" ("team_id")
    WHERE "assigned_admin_at" IS NOT NULL;
CREATE INDEX IF NOT EXISTS "trgm_idx_member_user_nickname"
    ON "t_member" USING gin ("user_nickname" gin_trgm_ops);

-- Create a default member for super admin in initial team.
INSERT INTO "t_member" (
    "id",
    "user_id",
    "user_nickname",
    "team_id",
    "assigned_raw_provider_at",
    "assigned_translator_at",
    "assigned_proofreader_at",
    "assigned_typesetter_at",
    "assigned_redrawer_at",
    "assigned_reviewer_at",
    "assigned_publisher_at",
    "assigned_admin_at"
) VALUES (
    'member-00000000-0000-0000-0000-000000000001',
    'user-00000000-0000-0000-0000-000000000001',
    'SuperAdmin-OvO',
    'team-00000000-0000-0000-0000-000000000001',
    NOW(),
    NOW(),
    NOW(),
    NOW(),
    NOW(),
    NOW(),
    NOW(),
    NOW()
) ON CONFLICT (id) DO NOTHING;

CREATE TABLE IF NOT EXISTS "t_member" (
    "id" TEXT PRIMARY KEY,
    
    "user_id" TEXT NOT NULL REFERENCES "t_user" ("id") ON DELETE CASCADE,
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


CREATE UNIQUE INDEX "uidx_member_user_team" ON "member_table" ("user_id", "team_id");

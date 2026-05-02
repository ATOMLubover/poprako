CREATE TABLE IF NOT EXISTS "t_team" (
    "id" TEXT PRIMARY KEY,
    
    "name" TEXT NOT NULL UNIQUE,
    "description" TEXT,

    "avatar_key" TEXT,
    "avatar_uploaded" BOOLEAN DEFAULT FALSE, 

    "workset_next_index" INTEGER NOT NULL DEFAULT 0,

    "created_at" TIMESTAMPTZ DEFAULT NOW(),
    "updated_at" TIMESTAMPTZ DEFAULT NOW()
);

-- Create a default team directly in database.
INSERT INTO "t_team" (
    "id",
    "name",
    "description"
) VALUES (
    'team-00000000-0000-0000-0000-000000000001',
    'PRTS 汉化组',
    '测测你的'
) ON CONFLICT (id) DO NOTHING;

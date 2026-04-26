CREATE TABLE IF NOT EXISTS "t_team" (
    "id" TEXT PRIMARY KEY,
    
    "name" VARCHAR(255) NOT NULL UNIQUE,
    "description" TEXT,

    "avatar_key" TEXT,
    "avatar_uploaded" BOOLEAN DEFAULT FALSE, 

    "created_at" TIMESTAMPTZ DEFAULT NOW(),
    "updated_at" TIMESTAMPTZ DEFAULT NOW()
);

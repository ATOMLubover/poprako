CREATE TABLE "invitation_table" (
    "id" TEXT PRIMARY KEY,

    -- 考虑到邀请不应该因为用户被删除而丢失，因此设置为 ON DELETE SET NULL
    "invitor_id" TEXT REFERENCES "user_table" ("id") ON DELETE SET NULL,
    "invitee_qq" TEXT NOT NULL,

    "invitation_code" TEXT UNIQUE NOT NULL,

    "to_be_picture_source" BOOLEAN NOT NULL DEFAULT FALSE,
    "to_be_translator" BOOLEAN NOT NULL DEFAULT FALSE,
    "to_be_proofreader" BOOLEAN NOT NULL DEFAULT FALSE,
    "to_be_typesetter" BOOLEAN NOT NULL DEFAULT FALSE,
    "to_be_reviewer" BOOLEAN NOT NULL DEFAULT FALSE,
    "to_be_admin" BOOLEAN NOT NULL DEFAULT FALSE,
    "to_be_super_admin" BOOLEAN NOT NULL DEFAULT FALSE,

    "pending" BOOLEAN NOT NULL DEFAULT TRUE,

    "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

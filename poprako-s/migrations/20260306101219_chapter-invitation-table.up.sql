CREATE TABLE "chapter_invitation" (
    id TEXT PRIMARY KEY,

    chapter_id TEXT NOT NULL REFERENCES "chapter_table" (id) ON DELETE CASCADE,

    inviter_id TEXT NOT NULL,
    invitee_qq TEXT NOT NULL,

    invitation_code TEXT NOT NULL,
    pending BOOLEAN NOT NULL DEFAULT TRUE,
    
    to_be_raw_provider BOOLEAN NOT NULL DEFAULT FALSE,
    to_be_translator BOOLEAN NOT NULL DEFAULT FALSE,
    to_be_proofreader BOOLEAN NOT NULL DEFAULT FALSE,
    to_be_typesetter BOOLEAN NOT NULL DEFAULT FALSE,
    to_be_redrawer BOOLEAN NOT NULL DEFAULT FALSE,
    to_be_reviewer BOOLEAN NOT NULL DEFAULT FALSE,
    to_be_publisher BOOLEAN NOT NULL DEFAULT FALSE,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX "idx_chapter_invitation_invitation_code"
    ON "chapter_invitation" (invitation_code);

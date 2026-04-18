CREATE TABLE "notification_table" (
    "id"                 TEXT        NOT NULL PRIMARY KEY,

    "user_id"            TEXT        NOT NULL REFERENCES "user_table" ("id") ON DELETE CASCADE,

    "content"            TEXT        NOT NULL,
    "is_read"            BOOLEAN     NOT NULL DEFAULT FALSE,

    "created_at"         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updated_at"         TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX "idx_notification_user_id_created_at_desc"
    ON "notification_table" ("user_id", "created_at" DESC);

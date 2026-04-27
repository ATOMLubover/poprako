CREATE TABLE IF NOT EXISTS "t_oss_message" (
    "id" TEXT PRIMARY KEY,

    "resource_type" TEXT NOT NULL,
    "resource_id" TEXT NOT NULL,
    "operation" TEXT NOT NULL,

    "status" TEXT NOT NULL,

    "object_keys" TEXT[] NOT NULL,

    "visible_at" TIMESTAMPTZ NOT NULL,
    "expire_at" TIMESTAMPTZ NOT NULL,
    "processing_at" TIMESTAMPTZ,

    "attempt_count" INTEGER NOT NULL DEFAULT 0,
    "last_error" TEXT,
  
    "created_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS "idx_oss_message_status_operation_visible" ON "t_oss_message" ("status", "operation", "visible_at");
CREATE INDEX IF NOT EXISTS "idx_oss_message_res_id_res_type_operation" ON "t_oss_message" ("resource_id", "resource_type", "operation");
CREATE INDEX IF NOT EXISTS "idx_oss_message_created_at" ON "t_oss_message" ("created_at");

-- Prevent multiple pending messages for the same resource and operation.
CREATE UNIQUE INDEX IF NOT EXISTS "uidx_oss_message_active_res_op"
    ON "t_oss_message" ("resource_type", "resource_id", "operation")
    WHERE "status" <> 'oss_messsage_status:completed';

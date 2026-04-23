CREATE TABLE oss_message_table (
    id                    VARCHAR(64)   NOT NULL PRIMARY KEY,

    resource_type         VARCHAR(64)   NOT NULL,
    resource_id           VARCHAR(64)   NOT NULL,
    operation             VARCHAR(32)   NOT NULL,

    status                VARCHAR(16)   NOT NULL DEFAULT 'pending',

    object_key            TEXT          NOT NULL DEFAULT '',
    payload_json          TEXT          NOT NULL DEFAULT '',

    visible_at            TIMESTAMPTZ(3) NOT NULL,
    expire_at             TIMESTAMPTZ(3),

    processing_at         TIMESTAMPTZ(3),
    attempt_count         INT           NOT NULL DEFAULT 0,

    last_error            TEXT          NOT NULL DEFAULT '',

    created_at            TIMESTAMPTZ(3) NOT NULL,
    updated_at            TIMESTAMPTZ(3) NOT NULL
);

CREATE INDEX idx_oss_message_status_visible ON oss_message_table (status, visible_at);
CREATE INDEX idx_oss_message_resource ON oss_message_table (resource_type, resource_id, operation);

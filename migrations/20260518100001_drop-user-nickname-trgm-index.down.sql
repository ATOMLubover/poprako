CREATE INDEX IF NOT EXISTS "trgm_idx_user_nickname"
    ON "t_user" USING gin ("nickname" gin_trgm_ops);

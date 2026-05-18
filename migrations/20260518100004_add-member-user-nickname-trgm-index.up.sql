CREATE INDEX IF NOT EXISTS "trgm_idx_member_user_nickname"
    ON "t_member" USING gin ("user_nickname" gin_trgm_ops);

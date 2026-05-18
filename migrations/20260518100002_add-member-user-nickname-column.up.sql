ALTER TABLE "t_member"
    ADD COLUMN IF NOT EXISTS "user_nickname" TEXT;

UPDATE "t_member" AS m
    SET "user_nickname" = u."nickname"
FROM "t_user" AS u
WHERE m."user_id" = u."id"
  AND m."user_nickname" IS NULL;

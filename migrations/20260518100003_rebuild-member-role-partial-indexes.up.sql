DROP INDEX IF EXISTS "idx_member_team_raw_provider";
DROP INDEX IF EXISTS "idx_member_team_translator";
DROP INDEX IF EXISTS "idx_member_team_proofreader";
DROP INDEX IF EXISTS "idx_member_team_typesetter";
DROP INDEX IF EXISTS "idx_member_team_redrawer";
DROP INDEX IF EXISTS "idx_member_team_reviewer";
DROP INDEX IF EXISTS "idx_member_team_publisher";
DROP INDEX IF EXISTS "idx_member_team_admin";

CREATE INDEX IF NOT EXISTS "idx_member_team_raw_provider"
    ON "t_member" ("team_id")
    WHERE "assigned_raw_provider_at" IS NOT NULL;
CREATE INDEX IF NOT EXISTS "idx_member_team_translator"
    ON "t_member" ("team_id")
    WHERE "assigned_translator_at" IS NOT NULL;
CREATE INDEX IF NOT EXISTS "idx_member_team_proofreader"
    ON "t_member" ("team_id")
    WHERE "assigned_proofreader_at" IS NOT NULL;
CREATE INDEX IF NOT EXISTS "idx_member_team_typesetter"
    ON "t_member" ("team_id")
    WHERE "assigned_typesetter_at" IS NOT NULL;
CREATE INDEX IF NOT EXISTS "idx_member_team_redrawer"
    ON "t_member" ("team_id")
    WHERE "assigned_redrawer_at" IS NOT NULL;
CREATE INDEX IF NOT EXISTS "idx_member_team_reviewer"
    ON "t_member" ("team_id")
    WHERE "assigned_reviewer_at" IS NOT NULL;
CREATE INDEX IF NOT EXISTS "idx_member_team_publisher"
    ON "t_member" ("team_id")
    WHERE "assigned_publisher_at" IS NOT NULL;
CREATE INDEX IF NOT EXISTS "idx_member_team_admin"
    ON "t_member" ("team_id")
    WHERE "assigned_admin_at" IS NOT NULL;

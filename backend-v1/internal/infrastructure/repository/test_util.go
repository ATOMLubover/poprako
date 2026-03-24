package repository

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"labelplus-next-web-be/internal/config"
	intf "labelplus-next-web-be/internal/domain/repository"

	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

var (
	testExecutor intf.Executor
	testDBMutex  sync.Mutex
	testDBOnce   sync.Once
	testDBErr    error
)

func ensureRepositoryTestDB(t *testing.T) intf.Executor {
	t.Helper()

	testDBOnce.Do(func() {
		ctx := context.Background()

		container, err := postgres.Run(
			ctx,
			"postgres:18-alpine",
			postgres.WithDatabase("test_db"),
			postgres.WithUsername("test_user"),
			postgres.WithPassword("test_password"),
			postgres.BasicWaitStrategies(),
		)
		if err != nil {
			testDBErr = fmt.Errorf("启动 Mock 数据库失败: %w", err)
			return
		}

		connStr, err := container.ConnectionString(ctx)
		if err != nil {
			testDBErr = fmt.Errorf("获取数据库连接字符串失败: %w", err)
			return
		}

		testExecutor, err = NewDatabaseExecutor(&config.DatabaseConfig{
			DatabaseURL:        connStr,
			MinIdleConnections: 1,
			MaxOpenConnections: 5,
		})
		if err != nil {
			testDBErr = fmt.Errorf("创建数据库执行器失败: %w", err)
			return
		}

		if err := prepareRepositoryTestSchema(testExecutor); err != nil {
			testDBErr = fmt.Errorf("初始化测试表结构失败: %w", err)
			return
		}
	})

	require.NoError(t, testDBErr)
	require.NotNil(t, testExecutor)

	return testExecutor
}

func prepareRepositoryTestSchema(executor intf.Executor) error {
	statements := []string{
		`CREATE EXTENSION IF NOT EXISTS pg_trgm`,
		`CREATE TABLE IF NOT EXISTS "user_table" (
			"id"                 TEXT        PRIMARY KEY,
			"name"               TEXT        NOT NULL,
			"qq"                 TEXT        UNIQUE NOT NULL,
			"avatar_oss_key"     TEXT        NOT NULL,
			"is_avatar_uploaded" BOOLEAN     NOT NULL DEFAULT FALSE,
			"password_hash"      TEXT        NOT NULL,
			"is_super_admin"     BOOLEAN     NOT NULL DEFAULT FALSE,
			"created_at"         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			"updated_at"         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			"deleted_at"         TIMESTAMPTZ
		)`,
		`CREATE TABLE IF NOT EXISTS "user_stats_table" (
			"id"                        TEXT        PRIMARY KEY,
			"user_id"                   TEXT        NOT NULL REFERENCES "user_table" ("id") ON DELETE CASCADE,
			"total_assignment_count"    INTEGER     NOT NULL DEFAULT 0,
			"active_assignment_count"   INTEGER     NOT NULL DEFAULT 0,
			"finished_assignment_count" INTEGER     NOT NULL DEFAULT 0,
			"created_at"                TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			"updated_at"                TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS "team_table" (
			"id"                 TEXT        PRIMARY KEY,
			"name"               TEXT        NOT NULL,
			"description"        TEXT,
			"avatar_oss_key"     TEXT,
			"is_avatar_uploaded" BOOLEAN     NOT NULL DEFAULT FALSE,
			"created_at"         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			"updated_at"         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			"deleted_at"         TIMESTAMPTZ
		)`,
		`CREATE TABLE IF NOT EXISTS "member_table" (
			"id"                       TEXT        PRIMARY KEY,
			"user_id"                  TEXT        NOT NULL REFERENCES "user_table" ("id") ON DELETE CASCADE,
			"team_id"                  TEXT        NOT NULL REFERENCES "team_table" ("id") ON DELETE CASCADE,
			"assigned_raw_provider_at" TIMESTAMPTZ,
			"assigned_translator_at"   TIMESTAMPTZ,
			"assigned_proofreader_at"  TIMESTAMPTZ,
			"assigned_typesetter_at"   TIMESTAMPTZ,
			"assigned_reviewer_at"     TIMESTAMPTZ,
			"assigned_publisher_at"    TIMESTAMPTZ,
			"assigned_admin_at"        TIMESTAMPTZ,
			"created_at"               TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			"updated_at"               TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			"deleted_at"               TIMESTAMPTZ
		)`,
		`CREATE TABLE IF NOT EXISTS "invitation_table" (
			"id"                  TEXT        PRIMARY KEY,
			"invitor_id"          TEXT        NOT NULL REFERENCES "user_table" ("id") ON DELETE CASCADE,
			"team_id"             TEXT        NOT NULL REFERENCES "team_table" ("id") ON DELETE CASCADE,
			"invitee_qq"          TEXT        NOT NULL,
			"invitation_code"     TEXT        UNIQUE NOT NULL,
			"to_be_raw_provider"  BOOLEAN     NOT NULL DEFAULT FALSE,
			"to_be_translator"    BOOLEAN     NOT NULL DEFAULT FALSE,
			"to_be_proofreader"   BOOLEAN     NOT NULL DEFAULT FALSE,
			"to_be_typesetter"    BOOLEAN     NOT NULL DEFAULT FALSE,
			"to_be_reviewer"      BOOLEAN     NOT NULL DEFAULT FALSE,
			"to_be_publisher"     BOOLEAN     NOT NULL DEFAULT FALSE,
			"to_be_admin"         BOOLEAN     NOT NULL DEFAULT FALSE,
			"pending"             BOOLEAN     NOT NULL DEFAULT TRUE,
			"created_at"          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			"updated_at"          TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS "workset_table" (
			"id"           TEXT        PRIMARY KEY,
			"team_id"      TEXT        NOT NULL REFERENCES "team_table" ("id") ON DELETE CASCADE,
			"index"        INTEGER     NOT NULL,
			"name"         TEXT        NOT NULL,
			"description"  TEXT,
			"comic_count"  INTEGER     NOT NULL DEFAULT 0,
			"created_at"   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			"updated_at"   TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,
		`CREATE UNIQUE INDEX IF NOT EXISTS "uidx_workset_team_id_index"
			ON "workset_table" ("team_id", "index")`,
		`CREATE TABLE IF NOT EXISTS "comic_table" (
			"id"                     TEXT        PRIMARY KEY,
			"workset_id"             TEXT        NOT NULL REFERENCES "workset_table" ("id") ON DELETE CASCADE,
			"index"                  INTEGER     NOT NULL DEFAULT 0,
			"title"                  TEXT        NOT NULL,
			"author"                 TEXT        NOT NULL,
			"composed_title"         TEXT        NOT NULL,
			"description"            TEXT,
			"chapter_count"          INTEGER     NOT NULL DEFAULT 0,
			"latest_uploaded_at"     TIMESTAMPTZ,
			"latest_transalating_at" TIMESTAMPTZ,
			"latest_translated_at"   TIMESTAMPTZ,
			"latest_proofreading_at" TIMESTAMPTZ,
			"latest_proofread_at"    TIMESTAMPTZ,
			"latest_typesetting_at"  TIMESTAMPTZ,
			"latest_typeset_at"      TIMESTAMPTZ,
			"latest_reviewed_at"     TIMESTAMPTZ,
			"latest_published_at"    TIMESTAMPTZ,
			"creator_id"             TEXT        NOT NULL REFERENCES "user_table" ("id") ON DELETE RESTRICT,
			"last_active_at"         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			"created_at"             TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			"updated_at"             TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			"deleted_at"             TIMESTAMPTZ
		)`,
		`CREATE UNIQUE INDEX IF NOT EXISTS "uidx_comic_workset_id_index"
			ON "comic_table" ("workset_id", "index")
			WHERE "deleted_at" IS NULL`,
		`CREATE TABLE IF NOT EXISTS "chapter_table" (
			"id"                    TEXT        PRIMARY KEY,
			"comic_id"              TEXT        NOT NULL REFERENCES "comic_table" ("id") ON DELETE CASCADE,
			"index"                 INTEGER     NOT NULL DEFAULT 0,
			"subtitle"              TEXT        NOT NULL,
			"page_count"            INTEGER     NOT NULL DEFAULT 0,
			"total_unit_count"      INTEGER     NOT NULL DEFAULT 0,
			"translated_unit_count" INTEGER     NOT NULL DEFAULT 0,
			"proofread_unit_count"  INTEGER     NOT NULL DEFAULT 0,
			"uploaded_at"           TIMESTAMPTZ,
			"transalating_at"       TIMESTAMPTZ,
			"translated_at"         TIMESTAMPTZ,
			"proofreading_at"       TIMESTAMPTZ,
			"proofread_at"          TIMESTAMPTZ,
			"typesetting_at"        TIMESTAMPTZ,
			"typeset_at"            TIMESTAMPTZ,
			"reviewed_at"           TIMESTAMPTZ,
			"published_at"          TIMESTAMPTZ,
			"creator_id"            TEXT        NOT NULL REFERENCES "user_table" ("id") ON DELETE RESTRICT,
			"created_at"            TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			"updated_at"            TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			"deleted_at"            TIMESTAMPTZ
		)`,
		`CREATE TABLE IF NOT EXISTS "page_table" (
			"id"                    TEXT        PRIMARY KEY,
			"chapter_id"            TEXT        NOT NULL REFERENCES "chapter_table" ("id") ON DELETE CASCADE,
			"index"                 INTEGER     NOT NULL DEFAULT 0,
			"oss_key"               TEXT,
			"uploaded"              BOOLEAN     NOT NULL DEFAULT FALSE,
			"total_unit_count"      INTEGER     NOT NULL DEFAULT 0,
			"translated_unit_count" INTEGER     NOT NULL DEFAULT 0,
			"proofread_unit_count"  INTEGER     NOT NULL DEFAULT 0,
			"creator_id"            TEXT        NOT NULL REFERENCES "user_table" ("id") ON DELETE RESTRICT,
			"created_at"            TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			"updated_at"            TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			UNIQUE ("chapter_id", "index")
		)`,
		`CREATE TABLE IF NOT EXISTS "assignment_table" (
			"id"                       TEXT        PRIMARY KEY,
			"chapter_id"               TEXT        NOT NULL REFERENCES "chapter_table" ("id") ON DELETE CASCADE,
			"user_id"                  TEXT        NOT NULL REFERENCES "user_table" ("id") ON DELETE CASCADE,
			"assigned_raw_provider_at" TIMESTAMPTZ,
			"assigned_translator_at"   TIMESTAMPTZ,
			"assigned_proofreader_at"  TIMESTAMPTZ,
			"assigned_typesetter_at"   TIMESTAMPTZ,
			"assigned_redrawer_at"     TIMESTAMPTZ,
			"assigned_reviewer_at"     TIMESTAMPTZ,
			"assigned_publisher_at"    TIMESTAMPTZ,
			"created_at"               TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			"updated_at"               TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			UNIQUE ("chapter_id", "user_id")
		)`,
		`CREATE TABLE IF NOT EXISTS "unit_table" (
			"id"                  TEXT        PRIMARY KEY,
			"page_id"             TEXT        NOT NULL REFERENCES "page_table" ("id") ON DELETE CASCADE,
			"x_coord"             REAL        NOT NULL,
			"y_coord"             REAL        NOT NULL,
			"index"               INTEGER     NOT NULL,
			"in_bubble"           BOOLEAN     NOT NULL DEFAULT TRUE,
			"is_proofread"        BOOLEAN     NOT NULL DEFAULT FALSE,
			"translated_text"     TEXT,
			"translator_id"       TEXT        REFERENCES "user_table" ("id") ON DELETE SET NULL,
			"translator_comment"  TEXT,
			"proofreader_text"    TEXT,
			"proofreader_id"      TEXT        REFERENCES "user_table" ("id") ON DELETE SET NULL,
			"proofreader_comment" TEXT,
			"created_at"          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			"updated_at"          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			UNIQUE ("page_id", "index")
		)`,
	}

	for _, statement := range statements {
		if err := executor.Exec(statement).Error; err != nil {
			return err
		}
	}

	return nil
}

func resetRepositoryTestTables(executor intf.Executor) error {
	statements := []string{
		`TRUNCATE TABLE "unit_table" RESTART IDENTITY CASCADE`,
		`TRUNCATE TABLE "page_table" RESTART IDENTITY CASCADE`,
		`TRUNCATE TABLE "assignment_table" RESTART IDENTITY CASCADE`,
		`TRUNCATE TABLE "chapter_table" RESTART IDENTITY CASCADE`,
		`TRUNCATE TABLE "comic_table" RESTART IDENTITY CASCADE`,
		`TRUNCATE TABLE "workset_table" RESTART IDENTITY CASCADE`,
		`TRUNCATE TABLE "invitation_table" RESTART IDENTITY CASCADE`,
		`TRUNCATE TABLE "member_table" RESTART IDENTITY CASCADE`,
		`TRUNCATE TABLE "team_table" RESTART IDENTITY CASCADE`,
		`TRUNCATE TABLE "user_stats_table" RESTART IDENTITY CASCADE`,
		`TRUNCATE TABLE "user_table" RESTART IDENTITY CASCADE`,
	}

	for _, statement := range statements {
		if err := executor.Exec(statement).Error; err != nil {
			return err
		}
	}

	return nil
}

func runRepositoryTest(t *testing.T, fn func(executor intf.Executor)) {
	t.Helper()

	executor := ensureRepositoryTestDB(t)

	testDBMutex.Lock()
	defer testDBMutex.Unlock()

	require.NoError(t, resetRepositoryTestTables(executor))

	fn(executor)

	if sqlDB, err := executor.DB(); err == nil {
		require.NoError(t, sqlDB.PingContext(context.Background()))
	}
}

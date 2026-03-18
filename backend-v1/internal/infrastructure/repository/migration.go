package repository

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	intf "labelplus-next-web-be/internal/domain/repository"
)

const migrationRecordTableName = "app_migration_record_table"

type migrationFile struct {
	Version string
	Path    string
	Name    string
}

// RunMigrations 自动执行 migrations 目录下的所有 up SQL 脚本。
//
// 规则：
// - 按文件名升序执行（与时间戳命名一致）
// - 每个脚本只执行一次，执行记录写入 app_migration_record_table
func RunMigrations(databaseExecutor intf.Executor) error {
	if err := ensureMigrationRecordTable(databaseExecutor); err != nil {
		return err
	}

	migrationFileList, err := loadMigrationFiles("migrations")
	if err != nil {
		return err
	}

	if len(migrationFileList) == 0 {
		return nil
	}

	appliedVersionSet, err := loadAppliedVersionSet(databaseExecutor)
	if err != nil {
		return err
	}

	if len(appliedVersionSet) == 0 {
		bootstrapped, bootstrapErr := bootstrapWhenSchemaAlreadyExists(databaseExecutor, migrationFileList)
		if bootstrapErr != nil {
			return bootstrapErr
		}
		if bootstrapped {
			return nil
		}
	}

	for _, migration := range migrationFileList {
		if _, exists := appliedVersionSet[migration.Version]; exists {
			continue
		}

		migrationSQLBytes, readErr := os.ReadFile(migration.Path)
		if readErr != nil {
			return fmt.Errorf("读取迁移脚本失败(%s): %w", migration.Name, readErr)
		}

		transactionErr := databaseExecutor.Transaction(func(transactionExecutor intf.Executor) error {
			if execErr := transactionExecutor.Exec(string(migrationSQLBytes)).Error; execErr != nil {
				return fmt.Errorf("执行迁移脚本失败(%s): %w", migration.Name, execErr)
			}

			now := time.Now()
			if insertErr := transactionExecutor.Exec(
				`INSERT INTO `+migrationRecordTableName+` (version, name, applied_at) VALUES (?, ?, ?)`,
				migration.Version,
				migration.Name,
				now,
			).Error; insertErr != nil {
				return fmt.Errorf("记录迁移执行结果失败(%s): %w", migration.Name, insertErr)
			}

			return nil
		})
		if transactionErr != nil {
			return transactionErr
		}
	}

	return nil
}

func ensureMigrationRecordTable(databaseExecutor intf.Executor) error {
	return databaseExecutor.Exec(`
CREATE TABLE IF NOT EXISTS app_migration_record_table (
	version    TEXT        PRIMARY KEY,
	name       TEXT        NOT NULL,
	applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
`).Error
}

func loadMigrationFiles(migrationDir string) ([]migrationFile, error) {
	migrationPathPattern := filepath.Join(migrationDir, "*.up.sql")
	pathList, err := filepath.Glob(migrationPathPattern)
	if err != nil {
		return nil, fmt.Errorf("扫描迁移脚本失败: %w", err)
	}

	result := make([]migrationFile, 0, len(pathList))
	for _, path := range pathList {
		name := filepath.Base(path)
		parts := strings.SplitN(name, "_", 2)
		if len(parts) != 2 || strings.TrimSpace(parts[0]) == "" {
			return nil, fmt.Errorf("迁移脚本命名不合法: %s", name)
		}

		result = append(result, migrationFile{
			Version: parts[0],
			Path:    path,
			Name:    name,
		})
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Name < result[j].Name
	})

	return result, nil
}

func loadAppliedVersionSet(databaseExecutor intf.Executor) (map[string]struct{}, error) {
	type row struct {
		Version string
	}

	var recordRowList []row
	if err := databaseExecutor.Raw(
		`SELECT version FROM ` + migrationRecordTableName + ` ORDER BY version ASC`,
	).Scan(&recordRowList).Error; err != nil {
		return nil, fmt.Errorf("读取迁移执行记录失败: %w", err)
	}

	result := make(map[string]struct{}, len(recordRowList))
	for _, item := range recordRowList {
		result[item.Version] = struct{}{}
	}

	return result, nil
}

func bootstrapWhenSchemaAlreadyExists(
	databaseExecutor intf.Executor,
	migrationFileList []migrationFile,
) (bool, error) {
	var exists bool
	if err := databaseExecutor.Raw(
		`SELECT to_regclass('public.user_table') IS NOT NULL`,
	).Scan(&exists).Error; err != nil {
		return false, fmt.Errorf("检查数据库现有结构失败: %w", err)
	}

	if !exists {
		return false, nil
	}

	transactionErr := databaseExecutor.Transaction(func(transactionExecutor intf.Executor) error {
		now := time.Now()
		for _, migration := range migrationFileList {
			if insertErr := transactionExecutor.Exec(
				`INSERT INTO `+migrationRecordTableName+` (version, name, applied_at) VALUES (?, ?, ?) ON CONFLICT (version) DO NOTHING`,
				migration.Version,
				migration.Name,
				now,
			).Error; insertErr != nil {
				return fmt.Errorf("初始化迁移记录失败(%s): %w", migration.Name, insertErr)
			}
		}

		return nil
	})
	if transactionErr != nil {
		return false, transactionErr
	}

	return true, nil
}

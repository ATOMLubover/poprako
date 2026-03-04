package repository

import (
	"labelplus-next-web-be/internal/config"
	intf "labelplus-next-web-be/internal/domain/repository"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDatabaseExecutor(databaseConfig *config.DatabaseConfig) (intf.Executor, error) {
	executor, err := gorm.Open(
		postgres.Open(databaseConfig.DatabaseURL),
		&gorm.Config{},
	)
	if err != nil {
		return nil, err
	}

	sqlDB, err := executor.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(databaseConfig.MinIdleConnections)
	sqlDB.SetMaxOpenConns(databaseConfig.MaxOpenConnections)

	return executor, nil
}

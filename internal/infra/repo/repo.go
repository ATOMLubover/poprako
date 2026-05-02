package repo_infra

import (
	"fmt"
	"time"

	"poprako-s/internal/cfg"

	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPgGdb(dbCfg *cfg.DbCfg) *gorm.DB {
	const PG_DB = "db_poprako_s"

	user := viper.GetString("DATABASE_USER")
	pwd := viper.GetString("DATABASE_PASSWORD")
	host := viper.GetString("DATABASE_HOST")
	port := viper.GetInt("DATABASE_PORT")

	dsn := fmt.Sprintf(
		"postgresql://%s:%s@%s:%d/%s?sslmode=disable",
		user, pwd, host, port, PG_DB,
	)

	gdb, err := gorm.Open(postgres.Open(dsn))
	if err != nil {
		panic(fmt.Sprintf("[NewPgGdb] failed to connect to postgres: %v", err))
	}

	return gdb
}

func ApplyCfg(gdb *gorm.DB, dbCfg *cfg.DbCfg) {
	sdb, err := gdb.DB()
	if err != nil {
		panic(fmt.Sprintf("[LoadCfg] failed to get sql.DB from gorm.DB: %v", err))
	}

	sdb.SetMaxOpenConns(dbCfg.MaxOpen)
	sdb.SetMaxIdleConns(dbCfg.MaxIdle)
	sdb.SetConnMaxLifetime(time.Second * time.Duration(dbCfg.MaxLife))
	sdb.SetConnMaxIdleTime(time.Second * time.Duration(dbCfg.MaxIdleTime))
}

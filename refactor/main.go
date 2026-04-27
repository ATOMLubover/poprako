// Package main
// @title Poprako-S Refactor API
// @version 0.4.0
// @description Poprako-S refactor API documentation
// @BasePath /api/v1
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
package main

import (
	"poprako-s/internal/api/http"
	"poprako-s/internal/api/state"
	app_impl "poprako-s/internal/app/impl"
	"poprako-s/internal/cfg"
	"poprako-s/internal/domain/svc"
	event_infra "poprako-s/internal/infra/event"
	oss_infra "poprako-s/internal/infra/ext/oss"
	token_infra "poprako-s/internal/infra/ext/token"
	repo_infra "poprako-s/internal/infra/repo"
	"poprako-s/internal/lgr"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {
	// Load environment variables from .env file if it exists,
	// otherwise rely on system environment variables.
	if err := godotenv.Load(); err != nil {
		zap.L().Warn(
			"[main] no .env file found, relying solely on environment variables",
		)
	}

	// Load application configuration.
	appCfg := cfg.NewAppCfg()
	if appCfg == nil {
		zap.L().Panic("[main] failed to load application configuration")
	}

	// Construct and set global logger based on application configuration.
	appLgr := lgr.New(appCfg)
	if appLgr == nil {
		zap.L().Panic("[main] failed to initialize logger")
	}

	lgr.SetGlobal(appLgr)

	// Load gorm database handle and apply config to it.
	gdb := repo_infra.NewPgGdb(appCfg.Db)
	if gdb == nil {
		zap.L().Panic("[main] failed to load gdb")
	}

	repo_infra.ApplyCfg(gdb, appCfg.Db)

	// Construct all repositories with `gdb`.
	txnCtrl := repo_infra.NewTxnCtrl(gdb)
	userRepo := repo_infra.NewUserRepo(gdb)
	teamRepo := repo_infra.NewTeamRepo(gdb)
	memberInvRepo := repo_infra.NewMemberInvRepo(gdb)
	memberRepo := repo_infra.NewMemberRepo(gdb)
	ossMsgRepo := repo_infra.NewOssMsgRepo(gdb)
	worksetRepo := repo_infra.NewWorksetRepo(gdb)
	// sysMailRepo := repo_infra.NewSysMailRepo(gdb)

	// Construct all services.
	userSvc := svc.NewUserSvc()
	memberSvc := svc.NewMemberSvc()
	ossMsgSvc := svc.NewOssMsgSvc()
	worksetSvc := svc.NewWorksetSvc()

	// Construct all external services.
	ossClient := oss_infra.NewOssClient()
	tknParser := token_infra.NewJwtParser()

	evBus := event_infra.NewEvBus()

	// Construct all applications with repositories and services.
	// DI deps are listed in order of repo, service, external service.
	userApp := app_impl.NewUserApp(
		txnCtrl, userRepo, memberRepo, memberInvRepo, ossMsgRepo,
		userSvc, memberSvc, ossMsgSvc,
		tknParser, ossClient, evBus,
	)
	teamApp := app_impl.NewTeamApp(
		teamRepo,
		ossClient,
	)
	worksetApp := app_impl.NewWorksetLogApp(
		app_impl.NewWorksetApp(txnCtrl, worksetSvc, memberRepo, worksetRepo),
	)

	appSt := state.NewAppState(
		appCfg,
		userApp,
		teamApp,
		worksetApp,
	)

	// Start HTTP server.
	http.RunServer(http.NewApp(appSt), appSt)
}

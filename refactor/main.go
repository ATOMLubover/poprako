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
	comicRepo := repo_infra.NewComicRepo(gdb)
	chapterRepo := repo_infra.NewChapterRepo(gdb)
	assignmentInvRepo := repo_infra.NewAssignmentInvRepo(gdb)
	assignmentRepo := repo_infra.NewAssignmentRepo(gdb)
	userStatsRepo := repo_infra.NewUserStatsRepo(gdb)
	sysMailRepo := repo_infra.NewSysMailRepo(gdb)

	// Construct all services.
	errClsf := repo_infra.NewErrClassifier()
	userSvc := svc.NewUserSvc()
	memberSvc := svc.NewMemberSvc()
	ossMsgSvc := svc.NewOssMsgSvc()
	worksetSvc := svc.NewWorksetSvc()
	comicSvc := svc.NewComicSvc()
	chapterSvc := svc.NewChapterSvc()
	assignmentInvSvc := svc.NewAssignmentInvSvc()
	assignmentSvc := svc.NewAssignmentSvc()

	// Construct all external services.
	ossClient := oss_infra.NewOssClient()
	tknParser := token_infra.NewJwtParser()

	// Initialize event bus and register handlers.
	// in case of unexpected shutdown, ensure event bus is closed to prevent resource leaks.
	evBus := event_infra.NewEvBus()

	defer evBus.Close()

	evBus.Sub(event_infra.NewNotifyInvitorHandler(teamRepo, sysMailRepo))
	evBus.Sub(event_infra.NewUpdateUserActiveHandler(userRepo))
	evBus.Sub(event_infra.NewUpdateUserStatsOnAssignmentCreatedHandler(userStatsRepo))
	evBus.Sub(event_infra.NewUpdateUserStatsOnAssignmentRemovedHandler(userStatsRepo))
	evBus.Sub(event_infra.NewUpdateUserStatsOnChapterPublishedHandler(userStatsRepo))
	evBus.Sub(event_infra.NewUpdateUserStatsOnChapterRemovedHandler(userStatsRepo))

	evBus.Run()

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
		app_impl.NewWorksetApp(txnCtrl, worksetSvc, memberRepo, worksetRepo, errClsf),
	)
	comicApp := app_impl.NewComicLogApp(
		app_impl.NewComicApp(txnCtrl, comicSvc, memberRepo, worksetRepo, comicRepo, errClsf),
	)
	chapterApp := app_impl.NewChapterLogApp(
		app_impl.NewChapterApp(txnCtrl, chapterSvc, assignmentSvc, memberRepo, worksetRepo, comicRepo, chapterRepo, assignmentRepo, evBus, errClsf),
	)
	userStatsApp := app_impl.NewUserStatsLogApp(
		app_impl.NewUserStatsApp(userStatsRepo),
	)
	sysMailApp := app_impl.NewSysMailLogApp(
		app_impl.NewSysMailApp(sysMailRepo),
	)
	assignmentInvApp := app_impl.NewAssignmentInvLogApp(
		app_impl.NewAssignmentInvApp(
			txnCtrl,
			assignmentInvSvc,
			assignmentSvc,
			userRepo,
			memberRepo,
			worksetRepo,
			comicRepo,
			chapterRepo,
			assignmentInvRepo,
			assignmentRepo,
			evBus,
			errClsf,
		),
	)
	assignmentApp := app_impl.NewAssignmentLogApp(
		app_impl.NewAssignmentApp(
			txnCtrl,
			assignmentSvc,
			memberRepo,
			worksetRepo,
			comicRepo,
			chapterRepo,
			assignmentRepo,
			evBus,
			errClsf,
		),
	)

	appSt := state.NewAppState(
		appCfg,
		userApp,
		teamApp,
		worksetApp,
		comicApp,
		chapterApp,
		userStatsApp,
		sysMailApp,
		assignmentInvApp,
		assignmentApp,
	)

	// Start HTTP server.
	http.RunServer(http.NewApp(appSt), appSt)
}

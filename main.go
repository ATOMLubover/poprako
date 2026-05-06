// Package main
// @title Poprako-S Refactor API
// @version 0.4.0
// @description Poprako-S refactor API documentation
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
package main

import (
	"context"
	"os/signal"
	"poprako-s/internal/api/http"
	"poprako-s/internal/api/state"
	app_impl "poprako-s/internal/app/impl"
	"poprako-s/internal/cfg"
	"poprako-s/internal/domain/svc"
	event_infra "poprako-s/internal/infra/event"
	oss_infra "poprako-s/internal/infra/ext/oss"
	token_infra "poprako-s/internal/infra/ext/token"
	repo_infra "poprako-s/internal/infra/repo"
	worker_infra "poprako-s/internal/infra/worker"
	"poprako-s/internal/lgr"
	"syscall"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {
	runCx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	defer stop()

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
	pageRepo := repo_infra.NewPageRepo(gdb)
	unitRepo := repo_infra.NewUnitRepo(gdb)
	assignmentInvRepo := repo_infra.NewAssignmentInvRepo(gdb)
	assignmentRepo := repo_infra.NewAssignmentRepo(gdb)
	sysMailRepo := repo_infra.NewSysMailRepo(gdb)

	// Construct all services.
	errClsf := repo_infra.NewErrClassifier()
	userSvc := svc.NewUserSvc()
	teamSvc := svc.NewTeamSvc()
	memberSvc := svc.NewMemberSvc()
	memberInvSvc := svc.NewMemberInvSvc()
	ossMsgSvc := svc.NewOssMsgSvc()
	worksetSvc := svc.NewWorksetSvc()
	comicSvc := svc.NewComicSvc()
	chapterSvc := svc.NewChapterSvc()
	pageSvc := svc.NewPageSvc()
	unitSvc := svc.UnitSvc{}
	assignmentInvSvc := svc.NewAssignmentInvSvc()
	assignmentSvc := svc.NewAssignmentSvc()
	exportSvc := svc.ChapterExportSvc{}
	importSvc := svc.ChapterImportSvc{}

	// Construct all external services.
	ossClient := oss_infra.NewOssClient()
	tknParser := token_infra.NewJwtParser()

	// Start OSS background worker for pending create/delete message consumption.
	go worker_infra.NewOssWorker(ossMsgRepo, ossClient).Run(runCx)

	// Initialize event bus and register handlers.
	// in case of unexpected shutdown, ensure event bus is closed to prevent resource leaks.
	evBus := event_infra.NewEvBus()

	defer evBus.Close()

	evBus.Sub(event_infra.NewNotifyInvitorHandler(teamRepo, sysMailRepo))
	evBus.Sub(event_infra.NewUpdateUserActiveHandler(userRepo))

	evBus.Run()

	// Construct all applications with repositories and services.
	// DI deps are listed in order of repo, service, external service.
	userApp := app_impl.NewUserLogApp(
		app_impl.NewUserApp(
			txnCtrl, userRepo, memberRepo, memberInvRepo, ossMsgRepo,
			userSvc, memberSvc, ossMsgSvc,
			tknParser, ossClient, evBus,
		),
	)
	teamApp := app_impl.NewTeamLogApp(
		app_impl.NewTeamApp(
			txnCtrl,
			teamRepo,
			userRepo,
			memberRepo,
			ossMsgRepo,
			teamSvc,
			ossMsgSvc,
			ossClient,
			errClsf,
		),
	)
	memberApp := app_impl.NewMemberLogApp(
		app_impl.NewMemberApp(
			txnCtrl,
			userRepo,
			memberRepo,
			memberInvRepo,
			memberSvc,
			ossClient,
			errClsf,
		),
	)
	memberInvApp := app_impl.NewMemberInvLogApp(
		app_impl.NewMemberInvApp(
			txnCtrl,
			memberRepo,
			memberInvRepo,
			memberInvSvc,
			errClsf,
		),
	)
	worksetApp := app_impl.NewWorksetLogApp(
		app_impl.NewWorksetApp(
			txnCtrl, memberRepo, worksetRepo,
			worksetSvc,
			evBus,
			errClsf,
		),
	)
	comicApp := app_impl.NewComicLogApp(
		app_impl.NewComicApp(
			txnCtrl, memberRepo, worksetRepo, comicRepo,
			ossClient,
			comicSvc, chapterSvc, assignmentSvc,
			evBus, errClsf,
		),
	)
	chapterApp := app_impl.NewChapterLogApp(
		app_impl.NewChapterApp(
			txnCtrl, memberRepo, worksetRepo, comicRepo, chapterRepo, assignmentRepo,
			chapterSvc, assignmentSvc,
			evBus, errClsf,
		),
	)
	chapterPortApp := app_impl.NewChapterPortLogApp(
		app_impl.NewChapterPortApp(
			txnCtrl, chapterRepo, comicRepo, pageRepo, unitRepo, assignmentRepo,
			unitSvc, exportSvc, importSvc,
			ossClient, errClsf,
		),
	)
	pageApp := app_impl.NewPageLogApp(
		app_impl.NewPageApp(
			txnCtrl, memberRepo, worksetRepo, comicRepo, chapterRepo, pageRepo, assignmentRepo,
			pageSvc, chapterSvc, ossMsgSvc,
			ossClient, errClsf,
		),
	)
	unitApp := app_impl.NewUnitLogApp(
		app_impl.NewUnitApp(
			txnCtrl, comicRepo, chapterRepo, pageRepo, unitRepo, assignmentRepo,
			unitSvc,
			errClsf,
		),
	)
	sysMailApp := app_impl.NewSysMailLogApp(
		app_impl.NewSysMailApp(sysMailRepo),
	)
	assignmentInvApp := app_impl.NewAssignmentInvLogApp(
		app_impl.NewAssignmentInvApp(
			txnCtrl, userRepo, memberRepo, worksetRepo, comicRepo, chapterRepo, assignmentInvRepo, assignmentRepo,
			assignmentInvSvc, assignmentSvc,
			evBus, errClsf,
		),
	)
	assignmentApp := app_impl.NewAssignmentLogApp(
		app_impl.NewAssignmentApp(
			txnCtrl, memberRepo, worksetRepo, comicRepo, chapterRepo, assignmentRepo,
			assignmentSvc,
			evBus, errClsf,
		),
	)

	appSt := state.NewAppState(
		appCfg,
		userApp,
		teamApp,
		memberApp,
		memberInvApp,
		worksetApp,
		comicApp,
		chapterApp,
		chapterPortApp,
		pageApp,
		unitApp,
		sysMailApp,
		assignmentInvApp,
		assignmentApp,
	)

	// Start HTTP server.
	http.RunServer(http.NewApp(appSt), appSt)
}

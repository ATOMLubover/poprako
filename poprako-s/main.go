// Package main
// @title PopRaKo-S API
// @version 0.0.1
// @description PopRaKo-S 后端 API 文档
// @BasePath /api/v1
//
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
package main

import (
	"fmt"

	api_http "poprako-s/internal/api/http"
	"poprako-s/internal/app"
	event_handler "poprako-s/internal/app/event_handler"
	"poprako-s/internal/cfg"
	"poprako-s/internal/domain/event"
	"poprako-s/internal/domain/service"
	event_infra "poprako-s/internal/infra/event"
	oss_infra "poprako-s/internal/infra/ext/oss"
	repo_infra "poprako-s/internal/infra/repo"
	"poprako-s/internal/lgr"
	app_state "poprako-s/internal/state"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {
	if err := godotenv.Load(); err != nil {
		zap.L().Warn("加载 .env 文件失败，可能是因为文件不存在")
	}

	appCfg, err := cfg.Load()
	if err != nil {
		panic("加载应用配置失败: " + err.Error())
	}

	lgr.Init(appCfg)

	// 初始化数据库连接
	gdb, err := repo_infra.NewGDB(appCfg.DB.DSN)
	if err != nil {
		panic(fmt.Sprintf("初始化数据库失败: %v", err))
	}

	// 初始化基础设施层：仓储
	userRepo := repo_infra.NewUserRepo(gdb)
	teamRepo := repo_infra.NewTeamRepo(gdb)
	memberRepo := repo_infra.NewMemberRepo(gdb)
	invRepo := repo_infra.NewMemberInvitationRepo(gdb)
	worksetRepo := repo_infra.NewWorksetRepo(gdb)
	comicRepo := repo_infra.NewComicRepo(gdb)
	chapterRepo := repo_infra.NewChapterRepo(gdb)
	pageRepo := repo_infra.NewPageRepo(gdb)
	assignmentRepo := repo_infra.NewAssignmentRepo(gdb)
	unitRepo := repo_infra.NewUnitRepo(gdb)
	chapterInvRepo := repo_infra.NewChapterInvitationRepo(gdb)
	txnMgr := repo_infra.NewTxnMgr(gdb)

	// 初始化 OSS 客户端（按 OSS_PLATFORM 选择实现）
	ossClient := oss_infra.NewClient()

	// 初始化事件总线
	eventBus, err := event_infra.NewEventBus()
	if err != nil {
		panic(fmt.Sprintf("初始化事件总线失败: %v", err))
	}
	defer eventBus.Close()

	// 注册事件处理器
	handlers := []event.EventHandler{
		event_handler.NewUnitSaveHandler(),
		event_handler.NewAssignmentCreateHandler(),
		event_handler.NewAssignmentRemoveHandler(),
		event_handler.NewComicCreateHandler(),
		event_handler.NewComicRemoveHandler(),
		event_handler.NewChapterCreateHandler(),
		event_handler.NewChapterCreatorAssignedHandler(),
		event_handler.NewChapterRemoveHandler(),
		event_handler.NewChapterPublishedHandler(),
	}
	for _, h := range handlers {
		if err := eventBus.SubUnsafe(h); err != nil {
			panic(fmt.Sprintf("注册事件处理器失败: %v", err))
		}
	}

	// 初始化领域服务
	userSvc := service.NewUserService()
	teamSvc := service.NewTeamService()
	memberSvc := service.NewMemberService()
	invSvc := service.NewMemberInvitationService()
	worksetSvc := service.NewWorksetService()
	comicSvc := service.NewComicService()
	chapterSvc := service.NewChapterService()
	chapterInvSvc := service.NewChapterInvitationService()
	pageSvc := service.NewPageService()
	assignmentSvc := service.NewAssignmentService()
	unitSvc := service.NewUnitService()

	// 初始化应用层
	userApp := app.NewUserApp(
		userSvc, memberSvc,
		userRepo, invRepo, memberRepo,
		txnMgr, eventBus, ossClient,
		&appCfg.Auth,
	)
	teamApp := app.NewTeamApp(
		teamSvc, memberSvc,
		userRepo, teamRepo, memberRepo,
		ossClient,
	)
	memberApp := app.NewMemberApp(
		memberSvc,
		userRepo, memberRepo, invRepo,
		txnMgr, ossClient,
	)
	invitationApp := app.NewInvitationApp(
		invSvc,
		memberRepo, invRepo,
	)
	worksetApp := app.NewWorksetApp(
		worksetSvc,
		memberRepo, worksetRepo,
		txnMgr,
	)
	comicApp := app.NewComicApp(
		comicSvc,
		memberRepo, worksetRepo, comicRepo,
		txnMgr, eventBus, ossClient,
	)
	chapterApp := app.NewChapterApp(
		chapterSvc, chapterInvSvc,
		memberRepo, worksetRepo, comicRepo, chapterRepo, assignmentRepo,
		userRepo, pageRepo, chapterInvRepo,
		txnMgr, eventBus, ossClient,
	)

	chapterExportApp := app.NewLogChapterExportApp(
		app.NewChapterExportApp(
			chapterRepo, comicRepo, pageRepo, unitRepo, assignmentRepo, ossClient,
		),
	)
	pageApp := app.NewPageApp(
		pageSvc,
		assignmentRepo, chapterRepo, pageRepo,
		ossClient,
	)
	assignmentApp := app.NewAssignmentApp(
		assignmentSvc,
		assignmentRepo, chapterInvRepo, chapterRepo, userRepo,
		txnMgr, eventBus, ossClient,
	)
	unitApp := app.NewUnitApp(
		unitSvc, eventBus,
		userRepo, pageRepo, chapterRepo, assignmentRepo, unitRepo,
		txnMgr,
	)

	// 组装 HTTP AppState
	state := app_state.NewAppState(
		appCfg,
		userApp, teamApp, memberApp, invitationApp,
		worksetApp, comicApp, chapterApp, chapterExportApp, pageApp,
		assignmentApp, unitApp,
	)

	zap.L().Info("应用状态初始化完成，HTTP 服务器启动")

	// 启动 HTTP 服务器（阻塞）
	if err := api_http.Serve(state); err != nil {
		panic(fmt.Sprintf("服务器停止: %v", err))
	}

	zap.L().Info("HTTP 服务器已正常退出")
}

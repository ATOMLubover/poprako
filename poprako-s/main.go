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
	app_state "poprako-s/internal/state"

	"go.uber.org/zap"
)

func main() {
	// 初始化全局 logger
	logger, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(logger)
	defer logger.Sync()

	// 加载应用配置
	appCfg, err := cfg.Load()
	if err != nil {
		panic(fmt.Sprintf("加载配置失败: %v", err))
	}

	// 初始化数据库连接
	db, err := repo_infra.NewGDB(appCfg.DB.DSN)
	if err != nil {
		panic(fmt.Sprintf("初始化数据库失败: %v", err))
	}

	// 初始化基础设施层：仓储
	userRepo := repo_infra.NewUserRepo(db)
	teamRepo := repo_infra.NewTeamRepo(db)
	memberRepo := repo_infra.NewMemberRepo(db)
	invRepo := repo_infra.NewInvitationRepo(db)
	worksetRepo := repo_infra.NewWorksetRepo(db)
	comicRepo := repo_infra.NewComicRepo(db)
	chapterRepo := repo_infra.NewChapterRepo(db)
	pageRepo := repo_infra.NewPageRepo(db)
	assignmentRepo := repo_infra.NewAssignmentRepo(db)
	unitRepo := repo_infra.NewUnitRepo(db)
	txnMgr := repo_infra.NewTxnMgr(db)

	// 初始化 OSS 客户端（当前为占位实现）
	ossClient := oss_infra.NewNoopClient()

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
	invSvc := service.NewInvitationService()
	worksetSvc := service.NewWorksetService()
	comicSvc := service.NewComicService()
	chapterSvc := service.NewChapterService()
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
		chapterSvc,
		memberRepo, worksetRepo, comicRepo, chapterRepo, assignmentRepo,
		userRepo, pageRepo,
		txnMgr, eventBus, ossClient,
	)
	pageApp := app.NewPageApp(
		pageSvc,
		assignmentRepo, chapterRepo, pageRepo,
		ossClient,
	)
	assignmentApp := app.NewAssignmentApp(
		assignmentSvc,
		assignmentRepo, chapterRepo, userRepo,
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
		worksetApp, comicApp, chapterApp, pageApp,
		assignmentApp, unitApp,
	)

	// 启动 HTTP 服务器（阻塞）
	if err := api_http.Serve(state); err != nil {
		panic(fmt.Sprintf("服务器停止: %v", err))
	}
}

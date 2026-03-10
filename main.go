// Package main
// @title LabelPlus Next Web API
// @version 0.1.0
// @description LabelPlus Next Web 后端 API 文档
// @BasePath /api/v1
//
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
package main

import (
	"labelplus-next-web-be/internal/api/http"
	"labelplus-next-web-be/internal/application"
	"labelplus-next-web-be/internal/config"
	repository_infra "labelplus-next-web-be/internal/infrastructure/repository"
	"labelplus-next-web-be/internal/log"
	"labelplus-next-web-be/internal/state"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {
	if err := godotenv.Load(); err != nil {
		panic("加载 .env 环境变量失败")
	}

	appConfig, err := config.Load()
	if err != nil {
		panic("加载应用配置失败: " + err.Error())
	}

	log.Init(appConfig)

	databaseExecutor, err := repository_infra.NewDatabaseExecutor(
		&appConfig.DatabaseConfig,
	)
	if err != nil {
		panic("创建数据库执行器失败: " + err.Error())
	}

	userRepository := repository_infra.NewUserRepository(databaseExecutor)
	teamRepository := repository_infra.NewTeamRepository(databaseExecutor)
	invitationRepository := repository_infra.NewInvitationRepository(databaseExecutor)
	memberRepository := repository_infra.NewMemberRepository(databaseExecutor)
	comicRepository := repository_infra.NewComicRepository(databaseExecutor)
	chapterRepository := repository_infra.NewChapterRepository(databaseExecutor)

	userApplication := application.NewUserApplication(
		&appConfig.AuthConfig,
		userRepository,
		memberRepository,
		invitationRepository,
	)
	teamApplication := application.NewTeamApplication(
		userRepository,
		teamRepository,
		memberRepository,
	)
	invitationApplication := application.NewInvitationApplication(
		userRepository,
		memberRepository,
		invitationRepository,
	)
	memberApplication := application.NewMemberApplication(
		userRepository,
		memberRepository,
		invitationRepository,
	)
	comicApplication := application.NewComicApplication(
		userRepository,
		memberRepository,
		comicRepository,
	)
	chapterApplication := application.NewChapterApplication(
		memberRepository,
		comicRepository,
		chapterRepository,
	)

	appState := state.NewAppState(
		appConfig,
		userApplication,
		invitationApplication,
		teamApplication,
		memberApplication,
		comicApplication,
		chapterApplication,
	)

	zap.L().Info("应用状态初始化完成，HTTP 服务器启动")

	if err := http.StartServer(appState); err != nil {
		panic("启动 HTTP 服务器失败: " + err.Error())
	}

	zap.L().Info("HTTP 服务器已正常退出")
}

package http

import (
	"labelplus-next-web-be/internal/config"
	"labelplus-next-web-be/internal/state"

	"github.com/iris-contrib/swagger/swaggerFiles"
	"github.com/iris-contrib/swagger/v12"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/recover"
	"github.com/kataras/iris/v12/middleware/requestid"

	_ "labelplus-next-web-be/docs"
)

func StartServer(appState *state.AppState) error {
	app := initialize(appState)

	if err := app.Listen(appState.AppConfig.ServerAddress); err != nil {
		return err
	}

	return nil
}

func initialize(appState *state.AppState) *iris.Application {
	app := iris.Default()

	// 启用 request ID 和 panic 恢复中间件
	app.Use(requestid.New())
	// 启用日志记录中间件
	app.Use(LogMiddleware(appState))
	// 启用 panic 恢复中间件
	app.Use(recover.New())
	// 启用错误记录中间件

	// 设置路由
	apiParty := app.Party("/api/v1")

	// 认证相关路由（无需登录）
	authParty := apiParty.Party("/auth")
	{
		authParty.Post("/login", Login(appState))
		authParty.Post("/register", Register(appState))
	}

	// 以下路由需要登录
	authorizedParty := apiParty.Party("/", AuthorizeMiddleware(appState))

	// 用户相关路由
	userParty := authorizedParty.Party("/users")
	{
		userParty.Get("/{user_id}", GetUserByID(appState))
		userParty.Get("/", ListUsers(appState))
		userParty.Delete("/{user_id}", RemoveUserByID(appState))
	}

	// 汉化组相关路由
	teamParty := authorizedParty.Party("/teams")
	{
		teamParty.Post("/", CreateTeam(appState))
		teamParty.Get("/", ListAllTeams(appState))
		teamParty.Get("/mine", ListMyTeams(appState))
		teamParty.Patch("/{team_id}", UpdateTeam(appState))
		teamParty.Delete("/{team_id}", DeleteTeam(appState))
	}

	// 成员相关路由
	memberParty := authorizedParty.Party("/members")
	{
		memberParty.Post("/", CreateMember(appState))
		memberParty.Get("/", ListMembers(appState))
		memberParty.Patch("/{member_id}", UpdateMemberRole(appState))
		memberParty.Delete("/{member_id}", RemoveMember(appState))
	}

	// 邀请相关路由
	invitationParty := authorizedParty.Party("/invitations")
	{
		invitationParty.Get("/", ListInvitations(appState))
		invitationParty.Post("/", CreateInvitation(appState))
		invitationParty.Patch("/{invitation_id}", PatchInvitation(appState))
		invitationParty.Delete("/{invitation_id}", DeleteInvitation(appState))
	}

	// 初始化 Swagger UI
	initializeSwagger(app, appState.AppConfig)

	return app
}

func initializeSwagger(app *iris.Application, appConfig *config.AppConfig) {
	if appConfig.IsProduction() {
		return
	}

	// 在非生产环境中启用 Swagger UI
	app.Get(
		"/swagger/{any:path}",
		swagger.WrapHandler(swaggerFiles.Handler, func(c *swagger.Config) {
			// 指定获取 Swagger 文档的 URL
			c.URL = "/swagger/doc.json" // Modified to relative path
		}),
	)
}

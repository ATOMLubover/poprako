package http

import (
	_ "poprako-s/docs"
	"poprako-s/internal/cfg"
	"poprako-s/internal/state"

	"github.com/iris-contrib/swagger/swaggerFiles"
	"github.com/iris-contrib/swagger/v12"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/recover"
	"github.com/kataras/iris/v12/middleware/requestid"
)

// Serve 启动 HTTP 服务器，阻塞直到服务器停止
func Serve(appState *state.AppState) error {
	app := initApp(appState)

	if err := app.Listen(appState.Cfg.HTTPAddr); err != nil {
		return err
	}

	return nil
}

func initApp(appState *state.AppState) *iris.Application {
	app := iris.Default()

	// 启用 request ID 和 panic 恢复中间件
	app.Use(requestid.New())
	// 启用日志记录中间件
	app.Use(LogMiddleware(appState))
	// 启用 panic 恢复中间件
	app.Use(recover.New())

	// 初始化 Swagger（非生产环境）
	initSwagger(app, appState.Cfg)

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
		// 兼容历史前端路径 /users/me
		userParty.Get("/me", GetMyUser(appState))
		// 兼容历史前端路径 /users/me/stats
		userParty.Get("/me/stats", GetMyUserStats(appState))
		userParty.Get("/mine", GetMyUser(appState))
		userParty.Get("/{user_id}", GetUserByID(appState))
		userParty.Get("/mine/stats", GetMyUserStats(appState))
		userParty.Post("/mine/avatar", ReserveMyAvatar(appState))
		userParty.Post("/mine/avatar/confirm", ConfirmMyAvatarUploaded(appState))
		userParty.Put("/mine", UpdateMyUser(appState))
		userParty.Delete("/{user_id}", RemoveUser(appState))
	}

	// 汉化组相关路由
	teamParty := authorizedParty.Party("/teams")
	{
		teamParty.Post("/", CreateTeam(appState))
		teamParty.Get("/", ListTeams(appState))
		teamParty.Get("/mine", ListMyTeams(appState))
		teamParty.Post("/{team_id}/avatar", ReserveTeamAvatar(appState))
		teamParty.Post("/{team_id}/avatar/confirm", ConfirmTeamAvatarUploaded(appState))
		teamParty.Put("/{team_id}", UpdateTeam(appState))
		teamParty.Delete("/{team_id}", DeleteTeam(appState))
	}

	// 成员相关路由
	memberParty := authorizedParty.Party("/members")
	{
		memberParty.Post("/", CreateMember(appState))
		memberParty.Post("/join", JoinTeam(appState))
		memberParty.Get("/mine", ListMyMembers(appState))
		memberParty.Get("/", ListMembers(appState))
		memberParty.Put("/{member_id}", UpdateMemberRole(appState))
		memberParty.Delete("/{member_id}", RemoveMember(appState))
	}

	// 邀请相关路由
	invitationParty := authorizedParty.Party("/invitations")
	{
		invitationParty.Get("/", ListInvitations(appState))
		invitationParty.Post("/", CreateInvitation(appState))
		invitationParty.Put("/{invitation_id}", PatchInvitation(appState))
		invitationParty.Delete("/{invitation_id}", DeleteInvitation(appState))
	}

	// 工作集相关路由
	worksetParty := authorizedParty.Party("/worksets")
	{
		worksetParty.Get("/", ListWorksets(appState))
		worksetParty.Post("/", CreateWorkset(appState))
		worksetParty.Put("/{workset_id}", UpdateWorkset(appState))
		worksetParty.Delete("/{workset_id}", DeleteWorkset(appState))
	}

	// 漫画相关路由
	comicParty := authorizedParty.Party("/comics")
	{
		comicParty.Get("/", ListComics(appState))
		comicParty.Post("/", CreateComic(appState))
		comicParty.Put("/{comic_id}", PatchComic(appState))
		comicParty.Delete("/{comic_id}", DeleteComic(appState))
		comicParty.Get("/{comic_id}/pinned-chapter", GetComicPinnedChapter(appState))
		comicParty.Post("/{comic_id}/cover", ReserveComicCover(appState))
		comicParty.Post("/{comic_id}/cover/confirm", ConfirmComicCoverUploaded(appState))
	}

	// 章节相关路由
	chapterParty := authorizedParty.Party("/chapters")
	{
		chapterParty.Get("/", ListComicChapters(appState))
		chapterParty.Post("/", CreateComicChapter(appState))
		chapterParty.Post("/{chapter_id}/invitations", InviteChapterAssignee(appState))
		chapterParty.Get("/{chapter_id}/export", ExportChapter(appState))
		chapterParty.Get("/{chapter_id}/export/lp", ExportChapterLp(appState))
		chapterParty.Post("/{chapter_id}/import", ImportChapter(appState))
		chapterParty.Patch("/{chapter_id}", UpdateChapter(appState))
		chapterParty.Delete("/{chapter_id}", DeleteComicChapter(appState))
	}

	// 页面相关路由
	pageParty := authorizedParty.Party("/pages")
	{
		pageParty.Get("/", ListChapterPages(appState))
		pageParty.Post("/", ReserveChapterPages(appState))
		pageParty.Put("/{page_id}", UpdatePage(appState))
		pageParty.Delete("/{page_id}", DeletePage(appState))
	}

	// 分配相关路由
	assignmentParty := authorizedParty.Party("/assignments")
	{
		assignmentParty.Get("/mine", ListMyAssignments(appState))
		assignmentParty.Get("/", ListChapterAssignments(appState))
		assignmentParty.Post("/join", JoinInvitorChapter(appState))
		assignmentParty.Post("/", CreateChapterAssignment(appState))
		assignmentParty.Delete("/{assignment_id}", RemoveAssignment(appState))
	}

	// unit 相关路由
	unitParty := authorizedParty.Party("/units")
	{
		unitParty.Get("/", ListPageUnits(appState))
		unitParty.Put("/", SavePageUnits(appState))
	}

	return app
}

func initSwagger(app *iris.Application, appCfg *cfg.AppCfg) {
	if appCfg.IsProduction() {
		return
	}

	app.Get(
		"/swagger/{any:path}",
		swagger.WrapHandler(swaggerFiles.Handler, func(c *swagger.Config) {
			c.URL = "/swagger/doc.json"
		}),
	)
}

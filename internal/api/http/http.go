package http

import (
	// _ "poprako-s/docs"
	"fmt"

	"poprako-s/internal/api/http/middleware"
	"poprako-s/internal/api/state"
	"poprako-s/internal/cfg"

	"github.com/iris-contrib/swagger/swaggerFiles"
	"github.com/iris-contrib/swagger/v12"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/recover"
	"github.com/kataras/iris/v12/middleware/requestid"

	_ "poprako-s/docs"
)

// `NewApp` builds the Iris HTTP application with all routes registered
// It wires middleware (`recover`, `requestid`, `LogLatency`, `Auth`), attaches
// swagger in dev mode, creates the `/api/v1` party tree, and returns
// the ready-to-serve `*iris.Application`
func NewApp(st *state.AppState) *iris.Application {
	app := iris.Default()

	// Enable(from first to last):
	// - panic recover
	// - request id
	// - log latency
	app.Use(recover.New())
	app.Use(requestid.New())
	app.Use(middleware.LogLatency(st.Cfg))

	// Enable CORS middleware.
	// app.UseRouter(middleware.CorsMiddleware())

	if st.Cfg.Env() == cfg.EnvDev {
		enableSwag(app)
	}

	apiV1 := app.Party("/api/v1")
	{
		auth := apiV1.Party("/auth")
		{
			auth.Post("/login", LoginUser(st))
			auth.Post("/register", RegUser(st))
		}

		// NOTE: all routes below require authorization.
		authorized := apiV1.Party("/", middleware.Auth())
		{
			user := authorized.Party("/users")
			{
				user.Get("/{user_id}", GetUserInfo(st))
				user.Get("/me", GetMyUserInfo(st))
				user.Put("/me", UpdateMyUserInfo(st))

				user.Post("/avatar", ResvUserAvatar(st))
				user.Post("/avatar/confirm", MarkUserAvatarUploaded(st))
			}

			team := authorized.Party("/teams")
			{
				team.Post("", CreateTeam(st))
				team.Get("", ListTeams(st))
				team.Get("/mine", ListMyTeams(st))
				team.Get("/{team_id}", GetTeamInfo(st))
				team.Put("/{team_id}", UpdateTeam(st))
				team.Post("/{team_id}/avatar", ReserveTeamAvatar(st))
				team.Post("/{team_id}/avatar/confirm", ConfirmTeamAvatarUploaded(st))
			}

			member := authorized.Party("/members")
			{
				member.Post("", CreateMember(st))
				member.Get("", ListTeamMembers(st))
				member.Get("/mine", ListMyMembers(st))
				member.Put("/{member_id}", UpdateMemberRole(st))
				member.Delete("/{member_id}", DeleteMember(st))
				member.Post("/join", JoinTeamByInvitation(st))
			}

			memberInvitation := authorized.Party("/member-invitations")
			{
				memberInvitation.Get("", ListMemberInvitations(st))
				memberInvitation.Post("", CreateMemberInvitation(st))
				memberInvitation.Put("/{invitation_id}", UpdateMemberInvitation(st))
				memberInvitation.Delete("/{invitation_id}", DeleteMemberInvitation(st))
			}

			workset := authorized.Party("/worksets")
			{
				workset.Get("", ListWorksets(st))
				workset.Post("", CreateWorkset(st))
				workset.Put("/{workset_id}", UpdateWorkset(st))
				workset.Delete("/{workset_id}", DeleteWorkset(st))
			}

			comic := authorized.Party("/comics")
			{
				comic.Get("", ListComics(st))
				comic.Get("/{comic_id}", GetComicById(st))
				comic.Post("", CreateComic(st))
				comic.Put("/{comic_id}", UpdateComic(st))
				comic.Post("/{comic_id}/cover", ResvComicCover(st))
				comic.Post("/{comic_id}/cover/confirm", MarkComicCoverUploaded(st))
				comic.Delete("/{comic_id}", DeleteComic(st))
			}

			chapter := authorized.Party("/chapters")
			{
				chapter.Get("", ListChapters(st))
				chapter.Get("/pinned", GetPinnedChapter(st))
				chapter.Get("/{chapter_id}/export", ExportChapter(st))
				chapter.Get("/{chapter_id}/export/lp", ExportChapterLp(st))
				chapter.Get("/{chapter_id}", GetChapterById(st))
				chapter.Post("/{chapter_id}/join", JoinChapter(st))
				chapter.Post("/{chapter_id}/import", ImportChapter(st))
				chapter.Post("", CreateChapter(st))
				chapter.Put("/{chapter_id}", UpdateChapter(st))
				chapter.Delete("/{chapter_id}", DeleteChapter(st))
			}

			page := authorized.Party("/pages")
			{
				page.Get("", ListChapterPages(st))
				page.Post("/reserve", ResvChapterPages(st))
				page.Post("/{page_id}/reserve", ResvChapterPage(st))
				page.Delete("", DeleteChapterPages(st))
				page.Post("/{page_id}/image/uploaded", MarkPageImageUploaded(st))
			}

			unit := authorized.Party("/units")
			{
				unit.Get("", ListPageUnits(st))
				unit.Post("", SavePageUnits(st))
			}

			sysMail := authorized.Party("/sys-mails")
			{
				sysMail.Get("", ListSysMail(st))
				sysMail.Post("/{sys_mail_id}/read", MarkSysMailRead(st))
			}

			assignmentInv := authorized.Party("/assignment-invitations")
			{
				assignmentInv.Get("", ListAssignmentInvitations(st))
				assignmentInv.Post("", CreateAssignmentInvitation(st))
				assignmentInv.Delete("/{invitation_id}", DeleteAssignmentInvitation(st))
				assignmentInv.Post("/join", JoinByAssignmentInvitation(st))
			}

			assignment := authorized.Party("/assignments")
			{
				assignment.Get("", ListAssignmentsByChapter(st))
				assignment.Get("/mine", ListMyAssignments(st))
				assignment.Put("", UpsertAssignment(st))
				assignment.Delete("/{assignment_id}", DeleteAssignment(st))
			}
		}
	}

	return app
}

// `RunServer` starts the Iris HTTP server on the configured host:port
// It panics if the server cannot start, making it suitable for `main.go`
func RunServer(app *iris.Application, st *state.AppState) {
	addr := fmt.Sprintf("%s:%d", st.Cfg.Http.Host, st.Cfg.Http.Port)

	if err := app.Run(iris.Addr(addr)); err != nil {
		panic(fmt.Sprintf("[RunServer] failed to start http server: %v", err))
	}
}

// `enableSwag` registers the swagger UI handler on the given app
// It is only called when the environment is `EnvDev`
func enableSwag(app *iris.Application) {
	app.Get(
		"/swagger/{any:path}",
		swagger.WrapHandler(
			swaggerFiles.Handler,
			func(c *swagger.Config) {
				c.URL = "/swagger/doc.json"
			},
		),
	)
}

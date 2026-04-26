package http

import (
	// _ "poprako-s/docs"
	"poprako-s/internal/api/http/middleware"
	"poprako-s/internal/api/state"
	"poprako-s/internal/cfg"

	"github.com/iris-contrib/swagger/swaggerFiles"
	"github.com/iris-contrib/swagger/v12"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/recover"
	"github.com/kataras/iris/v12/middleware/requestid"
)

func newApp(st *state.AppState) *iris.Application {
	app := iris.Default()

	// Enable(from first to last):
	// - panic recover
	// - request id
	// - log latency
	app.Use(recover.New())
	app.Use(requestid.New())
	app.Use(middleware.LogLatency(st.Cfg))

	if st.Cfg.Env == cfg.EnvDev {
		enableSwag(app)
	}

	apiV1 := app.Party("/api/v1")
	{
		auth := apiV1.Party("/auth")
		{
			auth.Post("/login", LoginUser(st))
			auth.Post("/reg", RegUser(st))
		}

		// NOTE: all routes below require authorization.
		authorized := apiV1.Party("/", middleware.Auth())
		{
			user := authorized.Party("/user")
			{
				user.Get("/{user_id}", GetUserInfo(st))
				user.Get("/me", GetMyUserInfo(st))

				user.Post("/avatar", ResvUserAvatar(st))
				user.Post("/avatar/confirm", MarkUserAvatarUploaded(st))
			}
		}
	}

	return app
}

func enableSwag(app *iris.Application) {
	app.Get(
		"/swagger/{any:path}",
		swagger.WrapHandler(swaggerFiles.Handler,
			func(c *swagger.Config) {
				c.URL = "/swagger/doc.json"
			}),
	)
}

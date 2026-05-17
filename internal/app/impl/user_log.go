package app_impl

import (
	"context"

	app_iface "poprako-s/internal/app"
	"poprako-s/internal/app/res"
	app_util "poprako-s/internal/app/util"
	"poprako-s/internal/app/val"

	"go.uber.org/zap"
)

// `userLogAppImpl` is a decorator for `UserApp` that logs the execution of its methods.
// It mainly logs args and protect inner from nil-problems.
type userLogAppImpl struct {
	inner app_iface.UserApp
}

// `NewUserLogApp` creates a logging decorator for `UserApp`.
func NewUserLogApp(inner app_iface.UserApp) app_iface.UserApp {
	// Ensure required dependency is provided.
	if inner == nil {
		zap.L().Panic("[NewUserLogApp] nil dependency", zap.Bool("inner", inner == nil))
	}

	return &userLogAppImpl{inner: inner}
}

// `GetInfo` logs request context and forwards the call to `inner`.
func (a *userLogAppImpl) GetInfo(cx context.Context, id string) app_res.AppRes[val.UserVal] {
	// Ensure `cx` is always non-nil for downstream calls.
	if cx == nil {
		cx = context.Background()
	}

	// Resolve logger from context and fallback to global logger.
	lgr := app_util.TakeLgr(cx)
	if lgr == nil {
		lgr = zap.L()
	}

	// Attach stable input fields and save logger back to context.
	lgr = lgr.With(zap.String("id", id))

	cx = app_util.SaveLgr(cx, lgr)

	return a.inner.GetInfo(cx, id)
}

// `Login` validates args, enriches logger context, and forwards the call.
func (a *userLogAppImpl) Login(cx context.Context, args *val.UserLoginArgs) app_res.AppRes[val.UserLoginRes] {
	// Ensure `cx` is always non-nil for downstream calls.
	if cx == nil {
		cx = context.Background()
	}

	// Reject nil args early to prevent nil dereference in `inner`.
	if args == nil {
		return app_res.Reject[val.UserLoginRes](app_res.BadRequest, "请求参数不能为空")
	}

	// Resolve logger from context and fallback to global logger.
	lgr := app_util.TakeLgr(cx)
	if lgr == nil {
		lgr = zap.L()
	}

	// Only log non-sensitive fields and save logger back to context.
	lgr = lgr.With(zap.String("args.qid", args.Qid))

	cx = app_util.SaveLgr(cx, lgr)

	return a.inner.Login(cx, args)
}

// `Register` validates args, enriches logger context, and forwards the call.
func (a *userLogAppImpl) Register(cx context.Context, args *val.UserRegArgs) app_res.AppRes[val.UserRegRes] {
	// Ensure `cx` is always non-nil for downstream calls.
	if cx == nil {
		cx = context.Background()
	}

	// Reject nil args early to prevent nil dereference in `inner`.
	if args == nil {
		return app_res.Reject[val.UserRegRes](app_res.BadRequest, "请求参数不能为空")
	}

	// Resolve logger from context and fallback to global logger.
	lgr := app_util.TakeLgr(cx)
	if lgr == nil {
		lgr = zap.L()
	}

	// Only log non-sensitive fields and save logger back to context.
	lgr = lgr.With(zap.String("args.qid", args.Qid))

	cx = app_util.SaveLgr(cx, lgr)

	return a.inner.Register(cx, args)
}

// `Update` validates args, enriches logger context, and forwards the call.
func (a *userLogAppImpl) Update(cx context.Context, args *val.UserUpdArgs) app_res.AppRes[app_res.None] {
	// Ensure `cx` is always non-nil for downstream calls.
	if cx == nil {
		cx = context.Background()
	}

	// Reject nil args early to prevent nil dereference in `inner`.
	if args == nil {
		return app_res.Reject[app_res.None](app_res.BadRequest, "请求参数不能为空")
	}

	// Resolve logger from context and fallback to global logger.
	lgr := app_util.TakeLgr(cx)
	if lgr == nil {
		lgr = zap.L()
	}

	// Attach stable input fields and save logger back to context.
	lgr = lgr.With(zap.String("args.id", args.Id), zap.String("args.qid", args.Qid))

	cx = app_util.SaveLgr(cx, lgr)

	return a.inner.Update(cx, args)
}

// `ResvAvatar` validates args, enriches logger context, and forwards the call.
func (a *userLogAppImpl) ResvAvatar(cx context.Context, currUid string, args *val.ResvUserAvatarArgs) app_res.AppRes[val.ResvUserAvatarRes] {
	// Ensure `cx` is always non-nil for downstream calls.
	if cx == nil {
		cx = context.Background()
	}

	// Reject nil args early to prevent nil dereference in `inner`.
	if args == nil {
		return app_res.Reject[val.ResvUserAvatarRes](app_res.BadRequest, "请求参数不能为空")
	}

	// Resolve logger from context and fallback to global logger.
	lgr := app_util.TakeLgr(cx)
	if lgr == nil {
		lgr = zap.L()
	}

	// Attach stable input fields and save logger back to context.
	lgr = lgr.With(zap.String("curr_uid", currUid), zap.String("args.user_id", args.UserId), zap.String("args.file_ext", args.FileExt))

	cx = app_util.SaveLgr(cx, lgr)

	return a.inner.ResvAvatar(cx, currUid, args)
}

// `TouchLastActive` enriches logger context and forwards the call.
func (a *userLogAppImpl) TouchLastActive(cx context.Context, id string) app_res.AppRes[app_res.None] {
	if cx == nil {
		cx = context.Background()
	}

	lgr := app_util.TakeLgr(cx)
	if lgr == nil {
		lgr = zap.L()
	}

	lgr = lgr.With(zap.String("id", id))

	cx = app_util.SaveLgr(cx, lgr)

	return a.inner.TouchLastActive(cx, id)
}

// `MarkAvatarUploaded` enriches logger context and forwards the call.
func (a *userLogAppImpl) MarkAvatarUploaded(cx context.Context, currUid string) app_res.AppRes[app_res.None] {
	// Ensure `cx` is always non-nil for downstream calls.
	if cx == nil {
		cx = context.Background()
	}

	// Resolve logger from context and fallback to global logger.
	lgr := app_util.TakeLgr(cx)
	if lgr == nil {
		lgr = zap.L()
	}

	// Attach stable input fields and save logger back to context.
	lgr = lgr.With(zap.String("curr_uid", currUid))

	cx = app_util.SaveLgr(cx, lgr)

	return a.inner.MarkAvatarUploaded(cx, currUid)
}

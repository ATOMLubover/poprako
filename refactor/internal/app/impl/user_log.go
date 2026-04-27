package app_impl

import (
	"context"

	app_iface "poprako-s/internal/app"
	"poprako-s/internal/app/res"
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
func (a *userLogAppImpl) GetInfo(cx context.Context, id string) res.AppRes[val.UserVal] {
	// Ensure `cx` is always non-nil for downstream calls.
	if cx == nil {
		cx = context.Background()
	}

	// Resolve logger from context and fallback to global logger.
	lgr := takeLgr(cx)
	if lgr == nil {
		lgr = zap.L()
	}

	// Attach stable input fields and save logger back to context.
	lgr = lgr.With(zap.String("id", id))
	cx = saveLgr(cx, lgr)

	return a.inner.GetInfo(cx, id)
}

// `Login` validates args, enriches logger context, and forwards the call.
func (a *userLogAppImpl) Login(cx context.Context, args *val.UserLoginArgs) res.AppRes[val.UserLoginRes] {
	// Ensure `cx` is always non-nil for downstream calls.
	if cx == nil {
		cx = context.Background()
	}

	// Reject nil args early to prevent nil dereference in `inner`.
	if args == nil {
		return res.Reject[val.UserLoginRes](res.ServerError, "登陆异常失败，请联系管理员")
	}

	// Resolve logger from context and fallback to global logger.
	lgr := takeLgr(cx)
	if lgr == nil {
		lgr = zap.L()
	}

	// Only log non-sensitive fields and save logger back to context.
	lgr = lgr.With(zap.String("args.qid", args.Qid))
	cx = saveLgr(cx, lgr)

	return a.inner.Login(cx, args)
}

// `Reg` validates args, enriches logger context, and forwards the call.
func (a *userLogAppImpl) Reg(cx context.Context, args *val.UserRegArgs) res.AppRes[val.UserRegRes] {
	// Ensure `cx` is always non-nil for downstream calls.
	if cx == nil {
		cx = context.Background()
	}

	// Reject nil args early to prevent nil dereference in `inner`.
	if args == nil {
		return res.Reject[val.UserRegRes](res.ServerError, "注册异常失败，请联系管理员")
	}

	// Resolve logger from context and fallback to global logger.
	lgr := takeLgr(cx)
	if lgr == nil {
		lgr = zap.L()
	}

	// Only log non-sensitive fields and save logger back to context.
	lgr = lgr.With(zap.String("args.qid", args.Qid))
	cx = saveLgr(cx, lgr)

	return a.inner.Reg(cx, args)
}

// `ResvAvatar` validates args, enriches logger context, and forwards the call.
func (a *userLogAppImpl) ResvAvatar(cx context.Context, args *val.ResvUserAvatarArgs) res.AppRes[val.ResvUserAvatarRes] {
	// Ensure `cx` is always non-nil for downstream calls.
	if cx == nil {
		cx = context.Background()
	}

	// Reject nil args early to prevent nil dereference in `inner`.
	if args == nil {
		return res.Reject[val.ResvUserAvatarRes](res.ServerError, "头像上传申请异常失败，请联系管理员")
	}

	// Resolve logger from context and fallback to global logger.
	lgr := takeLgr(cx)
	if lgr == nil {
		lgr = zap.L()
	}

	// Attach stable input fields and save logger back to context.
	lgr = lgr.With(zap.String("args.user_id", args.UserId), zap.String("args.file_ext", args.FileExt))
	cx = saveLgr(cx, lgr)

	return a.inner.ResvAvatar(cx, args)
}

// `MarkAvatarUploaded` enriches logger context and forwards the call.
func (a *userLogAppImpl) MarkAvatarUploaded(cx context.Context, currUid string) res.AppRes[res.None] {
	// Ensure `cx` is always non-nil for downstream calls.
	if cx == nil {
		cx = context.Background()
	}

	// Resolve logger from context and fallback to global logger.
	lgr := takeLgr(cx)
	if lgr == nil {
		lgr = zap.L()
	}

	// Attach stable input fields and save logger back to context.
	lgr = lgr.With(zap.String("curr_uid", currUid))
	cx = saveLgr(cx, lgr)

	return a.inner.MarkAvatarUploaded(cx, currUid)
}

package app_impl

import (
	"context"
	"fmt"

	app_iface "poprako-s/internal/app"
	"poprako-s/internal/app/res"
	app_util "poprako-s/internal/app/util"
	"poprako-s/internal/app/val"
	oss_iface "poprako-s/internal/domain/ext/oss"
	token_iface "poprako-s/internal/domain/ext/token"
	"poprako-s/internal/domain/model/aggr"
	"poprako-s/internal/domain/model/enum"
	repo_iface "poprako-s/internal/domain/repo"
	"poprako-s/internal/domain/svc"
	event_iface "poprako-s/internal/event"
	repo_infra "poprako-s/internal/infra/repo"

	"go.uber.org/zap"
)

type userAppImpl struct {
	txnCtrl       repo_iface.TxnCtrl
	userRepo      repo_iface.UserRepo
	memberRepo    repo_iface.MemberRepo
	memberInvRepo repo_iface.MemberInvRepo
	ossMsgRepo    repo_iface.OssMsgRepo

	userSvc   svc.UserSvc
	memberSvc svc.MemberSvc
	ossMsgSvc svc.OssMsgSvc

	tknParser token_iface.Parser
	ossSigner oss_iface.Signer

	evBus event_iface.EvBus
}

func NewUserApp(
	txnCtrl repo_iface.TxnCtrl,
	userRepo repo_iface.UserRepo,
	memberRepo repo_iface.MemberRepo,
	memberInvRepo repo_iface.MemberInvRepo,
	ossMsgRepo repo_iface.OssMsgRepo,
	userSvc svc.UserSvc,
	memberSvc svc.MemberSvc,
	ossMsgSvc svc.OssMsgSvc,
	tknParser token_iface.Parser,
	ossSigner oss_iface.Signer,
	evBus event_iface.EvBus,
) app_iface.UserApp {
	if txnCtrl == nil ||
		userRepo == nil ||
		memberRepo == nil ||
		memberInvRepo == nil ||
		ossMsgRepo == nil ||
		tknParser == nil ||
		evBus == nil ||
		ossSigner == nil {
		zap.L().Panic(
			"[NewUserApp] nil dependency",
			zap.Bool("txnCtrl", txnCtrl == nil),
			zap.Bool("userRepo", userRepo == nil),
			zap.Bool("memberRepo", memberRepo == nil),
			zap.Bool("memberInvRepo", memberInvRepo == nil),
			zap.Bool("ossMsgRepo", ossMsgRepo == nil),
			zap.Bool("tknParser", tknParser == nil),
			zap.Bool("evBus", evBus == nil),
			zap.Bool("ossSigner", ossSigner == nil),
		)
	}

	return &userAppImpl{
		txnCtrl:       txnCtrl,
		userRepo:      userRepo,
		memberRepo:    memberRepo,
		memberInvRepo: memberInvRepo,
		ossMsgRepo:    ossMsgRepo,
		userSvc:       userSvc,
		memberSvc:     memberSvc,
		ossMsgSvc:     ossMsgSvc,
		tknParser:     tknParser,
		evBus:         evBus,
		ossSigner:     ossSigner,
	}
}

// NOTE: no context except `err` should be attached to logger, as
// context should be recorded in logger wrapper impls.

func (a *userAppImpl) Login(cx context.Context, args *val.UserLoginArgs) res.AppRes[val.UserLoginRes] {
	lgr := app_util.TakeLgr(cx)

	creds, err := a.userRepo.GetCredsByQid(args.Qid)
	if err != nil {
		lgr.Error(
			"[userAppImpl.Login] failed to get user credentials",
			zap.Error(err),
		)
		return res.Reject[val.UserLoginRes](res.BadRequest, "用户不存在或密码错误")
	}

	if err := creds.VfyPwd(args.Pwd); err != nil {
		lgr.Error(
			"[userAppImpl.Login] failed to verify user password",
			zap.Error(err),
		)
		return res.Reject[val.UserLoginRes](res.BadRequest, "用户不存在或密码错误")
	}

	a.evBus.Pub(context.Background(), creds.PullEv())

	tk, err := a.tknParser.GenToken(&aggr.UserToken{
		UserId: creds.Id,
	})
	if err != nil {
		lgr.Error(
			"[userAppImpl.Login] failed to generate user token",
			zap.Error(err),
		)
		return res.Reject[val.UserLoginRes](res.ServerError, "生成用户令牌失败")
	}

	lgr.Info(
		"[userAppImpl.Login] user logged in",
		zap.String("user_id", creds.Id),
	)

	return res.Accept(&val.UserLoginRes{
		UserId: creds.Id,
		Token:  tk,
	})
}

func (a *userAppImpl) Reg(cx context.Context, args *val.UserRegArgs) res.AppRes[val.UserRegRes] {
	lgr := app_util.TakeLgr(cx)

	var (
		userId  string
		errCode = res.BadRequest
	)

	ev := make([]event_iface.Event, 0)

	if err := a.txnCtrl.RunWithTxn(func(cx context.Context) error {
		userRepo, err := repo_infra.TxnUserRepo(cx)
		if err != nil {
			errCode = res.ServerError
			return err
		}
		memberRepo, err := repo_infra.TxnMemberRepo(cx)
		if err != nil {
			errCode = res.ServerError
			return err
		}
		memberInvRepo, err := repo_infra.TxnMemberInvRepo(cx)
		if err != nil {
			errCode = res.ServerError
			return err
		}

		inv, err := memberInvRepo.GetPendingByInviteeQid(args.Qid)
		if err != nil {
			return err
		}

		if inv.InvCode != args.InvCode {
			return fmt.Errorf("wrong invitation code")
		}

		userReg, err := a.userSvc.NewUserReg(inv, args.Name, args.Pwd)
		if err != nil {
			return err
		}

		user, err := userRepo.Reg(userReg)
		if err != nil {
			return err
		}

		ev = append(ev, userReg.PullEv()...)

		memberCre := a.memberSvc.NewMemberCre(user.Id, inv.TeamId, inv.RoleMask)

		_, err = memberRepo.Create(memberCre)
		if err != nil {
			return err
		}

		if err = memberInvRepo.MarkCmpl(inv.Id); err != nil {
			return err
		}

		userId = user.Id

		return nil
	}); err != nil {
		lgr.Error(
			"[userAppImpl.Reg] failed to run registration transaction",
			zap.Error(err),
		)

		switch errCode {
		case res.ServerError:
			return res.Reject[val.UserRegRes](res.ServerError, "注册失败")
		default:
			return res.Reject[val.UserRegRes](res.BadRequest, err.Error())
		}
	}

	// If successfully registered, publish all events after transaction is committed.
	a.evBus.Pub(context.Background(), ev)

	tk, err := a.tknParser.GenToken(&aggr.UserToken{
		UserId: userId,
	})
	if err != nil {
		lgr.Error(
			"[userAppImpl.Reg] failed to generate user token",
			zap.Error(err),
		)
		return res.Reject[val.UserRegRes](res.ServerError, "生成用户令牌失败")
	}

	lgr.Info(
		"[userAppImpl.Reg] user registered",
		zap.String("user_id", userId),
	)

	return res.Accept(&val.UserRegRes{
		UserId: userId,
		Token:  tk,
	})
}

func (a *userAppImpl) GetInfo(cx context.Context, id string) res.AppRes[val.UserVal] {
	lgr := app_util.TakeLgr(cx)

	user, err := a.userRepo.GetById(id)
	if err != nil {
		lgr.Error(
			"[userAppImpl.GetInfo] failed to get user by id",
			zap.String("user_id", id),
			zap.Error(err),
		)
		return res.Reject[val.UserVal](res.BadRequest, "用户不存在")
	}

	userVal, err := asmUserVal(user, a.ossSigner)
	if err != nil {
		lgr.Error(
			"[userAppImpl.GetInfo] failed to assemble user info",
			zap.String("user_id", id),
			zap.Error(err),
		)
		return res.Reject[val.UserVal](res.ServerError, "获取用户信息失败")
	}

	return res.Accept(userVal)
}

// func (a *userAppImpl) UpdateInfo(cx context.Context, args *val.UserUpdateArgs) res.AppRes[res.None] {
// }

func (a *userAppImpl) ResvAvatar(cx context.Context, args *val.ResvUserAvatarArgs) res.AppRes[val.ResvUserAvatarRes] {
	lgr := app_util.TakeLgr(cx)

	key := fmt.Sprintf("user_avatar/%s.%s", args.UserId, args.FileExt)

	if err := a.txnCtrl.RunWithTxn(func(cx context.Context) error {
		userRepo, err := repo_infra.TxnUserRepo(cx)
		if err != nil {
			return err
		}
		ossMsgRepo, err := repo_infra.TxnOssMsgRepo(cx)
		if err != nil {
			return err
		}

		if err = userRepo.PrefillAvatarKey(args.UserId, key); err != nil {
			return err
		}

		if err = a.ossMsgSvc.SavePendingCre(ossMsgRepo, enum.OssResUserAvatar, args.UserId, []string{key}); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return res.Reject[val.ResvUserAvatarRes](res.ServerError, "生成头像上传信息失败")
	}

	// NOTE: as a pending creation message will be recycled after a certain period of time,
	// it's not a problem if the generated URL is not used within a short period of time,
	// as long as the URL is generated successfully and the corresponding pending
	// creation message is created in the database.
	url, err := a.ossSigner.GenPutUrl(key)
	if err != nil {
		lgr.Error(
			"[userAppImpl.ResvAvatar] failed to generate avatar upload URL",
			zap.Error(err),
		)

		return res.Reject[val.ResvUserAvatarRes](res.ServerError, "生成头像上传信息失败")
	}

	return res.Accept(&val.ResvUserAvatarRes{
		PutUrl: url,
	})
}

func (a *userAppImpl) MarkAvatarUploaded(cx context.Context, uid string) res.AppRes[res.None] {
	lgr := app_util.TakeLgr(cx)

	if err := a.txnCtrl.RunWithTxn(func(cx context.Context) error {
		userRepo, err := repo_infra.TxnUserRepo(cx)
		if err != nil {
			return err
		}
		ossMsgRepo, err := repo_infra.TxnOssMsgRepo(cx)
		if err != nil {
			return err
		}

		if err = userRepo.MarkAvatarUploaded(uid); err != nil {
			return err
		}

		if err = ossMsgRepo.MarkCmplByRes(enum.OssResUserAvatar, uid); err != nil {
			return err
		}

		return nil
	}); err != nil {
		lgr.Error(
			"[userAppImpl.MarkAvatarUploaded] failed to run transaction",
			zap.Error(err),
		)
		return res.Reject[res.None](res.ServerError, "标记头像上传状态失败")
	}

	return res.Accept(&res.None{})
}

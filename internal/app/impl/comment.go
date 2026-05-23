package app_impl

import (
	"context"

	app_iface "poprako-s/internal/app"
	app_res "poprako-s/internal/app/res"
	app_util "poprako-s/internal/app/util"
	"poprako-s/internal/app/val"
	"poprako-s/internal/domain/model/enum"
	"poprako-s/internal/domain/model/query"
	repo_iface "poprako-s/internal/domain/repo"
	"poprako-s/internal/domain/svc"

	"go.uber.org/zap"
)

// `commentAppImpl` is the default implementation of `CommentApp`.
type commentAppImpl struct {
	txnCtrl repo_iface.TxnCtrl

	commentSvc svc.CommentSvc

	memberRepo  repo_iface.MemberRepo
	commentRepo repo_iface.CommentRepo

	errClsf repo_iface.ErrClsf
}

// `NewCommentApp` creates one `CommentApp` implementation.
func NewCommentApp(
	txnCtrl repo_iface.TxnCtrl,
	memberRepo repo_iface.MemberRepo,
	commentRepo repo_iface.CommentRepo,
	commentSvc svc.CommentSvc,
	errClsf repo_iface.ErrClsf,
) app_iface.CommentApp {
	if txnCtrl == nil || memberRepo == nil || commentRepo == nil || errClsf == nil {
		zap.L().Panic(
			"[NewCommentApp] nil dependency",
			zap.Bool("txnCtrl", txnCtrl == nil),
			zap.Bool("memberRepo", memberRepo == nil),
			zap.Bool("commentRepo", commentRepo == nil),
			zap.Bool("errClsf", errClsf == nil),
		)
	}

	return &commentAppImpl{
		txnCtrl:     txnCtrl,
		commentSvc:  commentSvc,
		memberRepo:  memberRepo,
		commentRepo: commentRepo,
		errClsf:     errClsf,
	}
}

// `List` lists team comments.
func (a *commentAppImpl) List(cx context.Context, currUid string, args *val.ListCommentArgs) app_res.AppRes[[]val.CommentVal] {
	lgr := app_util.TakeLgr(cx)

	if re := vfyListCommentArgs(args); re.IsReject() {
		return app_res.Reject[[]val.CommentVal](re.Code(), re.Msg())
	}

	if re := a.commentSvc.CanListComment(currUid, args.TeamId, a.memberRepo, a.errClsf); re.IsReject() {
		return app_res.Reject[[]val.CommentVal](app_res.ErrCode(re.Code()), re.Msg())
	}

	items, err := a.commentRepo.List(
		query.ListCommentOpt{
			TeamId: args.TeamId,
			Pagi: query.PagiOpt{
				Offset: args.Offset,
				Limit:  args.Limit,
			},
		},
		enum.CommentInclUser,
	)
	if err != nil {
		lgr.Error("[commentAppImpl.List] failed to list comments", zap.Error(err))

		return app_res.Reject[[]val.CommentVal](app_res.ServerError, "获取留言列表失败")
	}

	vals := make([]val.CommentVal, len(items))
	for i := range items {
		vals[i] = asmCommentVal(&items[i])
	}

	return app_res.Accept(&vals)
}

// `Create` creates one team board comment.
func (a *commentAppImpl) Create(cx context.Context, currUid string, args *val.CreateCommentArgs) app_res.AppRes[val.CommentCreatedRes] {
	lgr := app_util.TakeLgr(cx)

	if re := vfyCreateCommentArgs(args); re.IsReject() {
		return app_res.Reject[val.CommentCreatedRes](re.Code(), re.Msg())
	}

	re, err := repo_iface.RunWithTxn[app_res.AppRes[val.CommentCreatedRes]](a.txnCtrl, func(prov repo_iface.Prov) (app_res.AppRes[val.CommentCreatedRes], error) {
		memberRepo := prov.MemberRepo()
		commentRepo := prov.CommentRepo()

		if re := a.commentSvc.CanCreateComment(currUid, args.TeamId, memberRepo, a.errClsf); re.IsReject() {
			return app_res.Reject[val.CommentCreatedRes](app_res.ErrCode(re.Code()), re.Msg()), app_res.DefErr()
		}

		commentCre, err := a.commentSvc.NewCommentCre(currUid, args.TeamId, args.Content)
		if err != nil {
			return app_res.Reject[val.CommentCreatedRes](app_res.BadRequest, err.Error()), app_res.DefErr()
		}

		comment, err := commentRepo.Create(commentCre)
		if err != nil {
			return app_res.Reject[val.CommentCreatedRes](app_res.ServerError, "发布留言失败"), err
		}

		return app_res.Accept(&val.CommentCreatedRes{Id: comment.Id}), nil
	})
	if err != nil {
		lgr.Error("[commentAppImpl.Create] failed to run transaction", zap.Error(err))

		return re
	}

	return re
}

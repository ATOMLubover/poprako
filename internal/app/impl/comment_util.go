package app_impl

import (
	app_res "poprako-s/internal/app/res"
	app_util "poprako-s/internal/app/util"
	"poprako-s/internal/app/val"
	"poprako-s/internal/domain/model/aggr"
)

// `asmCommentVal` converts one comment aggregate into app value object.
func asmCommentVal(item *aggr.Comment) val.CommentVal {
	var userVal *val.UserVal
	if item.User != nil {
		userVal = &val.UserVal{
			Id:             item.User.Id,
			Qid:            item.User.Qid,
			Nickname:       item.User.Nickname,
			AvatarUploaded: item.User.AvatarUploaded,
			IsSuperAdmin:   item.User.IsSuperAdmin,
			LastActiveAt:   item.User.LastActiveAt.UnixMilli(),
			CreatedAt:      item.User.CreatedAt.UnixMilli(),
			UpdatedAt:      item.User.UpdatedAt.UnixMilli(),
		}
	}

	return val.CommentVal{
		Id: item.Id,

		TeamId: item.TeamId,
		UserId: item.UserId,
		User:   userVal,

		Content: item.Content,

		CreatedAt: item.CreatedAt.UnixMilli(),
	}
}

// `vfyListCommentArgs` validates and normalizes comment list args.
func vfyListCommentArgs(args *val.ListCommentArgs) app_res.AppRes[app_res.None] {
	if args == nil || args.TeamId == "" {
		return app_res.Reject[app_res.None](app_res.BadRequest, "team_id 不能为空")
	}

	if re := app_util.ClampOffsetLimit(args.Offset, &args.Limit); re.IsReject() {
		return re
	}

	return app_res.Accept(&app_res.None{})
}

// `vfyCreateCommentArgs` validates create comment args.
func vfyCreateCommentArgs(args *val.CreateCommentArgs) app_res.AppRes[app_res.None] {
	if args == nil || args.TeamId == "" {
		return app_res.Reject[app_res.None](app_res.BadRequest, "team_id 不能为空")
	}

	if args.Content == "" {
		return app_res.Reject[app_res.None](app_res.BadRequest, "content 不能为空")
	}

	return app_res.Accept(&app_res.None{})
}

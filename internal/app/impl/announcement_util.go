package app_impl

import (
	app_res "poprako-s/internal/app/res"
	app_util "poprako-s/internal/app/util"
	"poprako-s/internal/app/val"
	"poprako-s/internal/domain/model/aggr"
)

// `asmAnnouncementVal` converts one announcement aggregate into app value object.
func asmAnnouncementVal(item *aggr.Announcement) val.AnnouncementVal {
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

	return val.AnnouncementVal{
		Id: item.Id,

		TeamId: item.TeamId,
		UserId: item.UserId,
		User:   userVal,

		Title:   item.Title,
		Content: item.Content,

		CreatedAt: item.CreatedAt.UnixMilli(),
	}
}

// `vfyListAnnouncementArgs` validates and normalizes announcement list args.
func vfyListAnnouncementArgs(args *val.ListAnnouncementArgs) app_res.AppRes[app_res.None] {
	if args == nil || args.TeamId == "" {
		return app_res.Reject[app_res.None](app_res.BadRequest, "team_id 不能为空")
	}

	if re := app_util.ClampOffsetLimit(args.Offset, &args.Limit); re.IsReject() {
		return re
	}

	return app_res.Accept(&app_res.None{})
}

// `vfyCreateAnnouncementArgs` validates create announcement args.
func vfyCreateAnnouncementArgs(args *val.CreateAnnouncementArgs) app_res.AppRes[app_res.None] {
	if args == nil || args.TeamId == "" {
		return app_res.Reject[app_res.None](app_res.BadRequest, "team_id 不能为空")
	}

	if args.Title == "" {
		return app_res.Reject[app_res.None](app_res.BadRequest, "title 不能为空")
	}

	if args.Content == "" {
		return app_res.Reject[app_res.None](app_res.BadRequest, "content 不能为空")
	}

	return app_res.Accept(&app_res.None{})
}

package app_impl

import (
	app_res "poprako-s/internal/app/res"
	app_util "poprako-s/internal/app/util"
	"poprako-s/internal/app/val"
	"poprako-s/internal/domain/model/aggr"
)

// `vfyListWorksetArgs` validates and normalizes list arguments.
func vfyListWorksetArgs(args *val.ListWorksetArgs) app_res.AppRes[app_res.None] {
	if args == nil {
		return app_res.Reject[app_res.None](app_res.BadRequest, "分页参数不能为空")
	}

	if re := app_util.ClampOffsetLimit(args.Offset, &args.Limit); re.IsReject() {
		return re
	}

	return app_res.Accept(&app_res.None{})
}

// `asmWorksetVal` converts a `Workset` aggregate into a `WorksetVal` value object.
func asmWorksetVal(ws *aggr.Workset) val.WorksetVal {
	var teamVal *val.TeamVal

	if ws.Team != nil {
		teamVal = &val.TeamVal{
			Id:             ws.Team.Id,
			Name:           ws.Team.Name,
			Desc:           ws.Team.Desc,
			AvatarUploaded: ws.Team.AvatarUploaded,
			CreatedAt:      ws.Team.CreatedAt.UnixMilli(),
			UpdatedAt:      ws.Team.UpdatedAt.UnixMilli(),
		}
	}

	return val.WorksetVal{
		Id:         ws.Id,
		TeamId:     ws.TeamId,
		Team:       teamVal,
		Index:      ws.Index,
		Name:       ws.Name,
		Desc:       ws.Desc,
		ComicCount: ws.ComicCount,
		CreatedAt:  ws.CreatedAt.UnixMilli(),
		UpdatedAt:  ws.UpdatedAt.UnixMilli(),
	}
}

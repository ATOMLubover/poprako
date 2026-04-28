package app_impl

import (
	"poprako-s/internal/app/res"
	app_util "poprako-s/internal/app/util"
	"poprako-s/internal/app/val"
	"poprako-s/internal/domain/model/aggr"
)


// `vfyListWorksetArgs` validates and normalizes list arguments.
func vfyListWorksetArgs(args *val.ListWorksetArgs) (res.ErrCode, string, bool) {
	if args == nil {
		return res.BadRequest, "分页参数不能为空", true
	}

	if code, msg, reject := app_util.ClampOffsetLimit(args.Offset, &args.Limit); reject {
		return code, msg, true
	}

	return 0, "", false
}

// `asmWorksetVal` converts a `Workset` aggregate into a `WorksetVal` value object.
func asmWorksetVal(ws *aggr.Workset) val.WorksetVal {
	return val.WorksetVal{
		Id:         ws.Id,
		TeamId:     ws.TeamId,
		Index:      ws.Index,
		Name:       ws.Name,
		Desc:       ws.Desc,
		ComicCount: ws.ComicCount,
		CreatedAt:  ws.CreatedAt.UnixMilli(),
		UpdatedAt:  ws.UpdatedAt.UnixMilli(),
	}
}

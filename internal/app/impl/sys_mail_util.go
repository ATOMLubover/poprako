package app_impl

import (
	app_res "poprako-s/internal/app/res"
	app_util "poprako-s/internal/app/util"
	"poprako-s/internal/app/val"
	"poprako-s/internal/domain/model/aggr"
)

// `asmSysMailVal` converts `SysMail` aggregate to app-facing value object.
func asmSysMailVal(mail *aggr.SysMail) val.SysMailVal {
	return val.SysMailVal{
		Id:        mail.Id,
		Title:     mail.Title,
		Content:   mail.Content,
		Read:      mail.Read,
		CreatedAt: mail.CreatedAt.UnixMilli(),
	}
}

// `vfyListSysMailArgs` validates and normalizes list arguments.
func vfyListSysMailArgs(args *val.ListSysMailArgs) app_res.AppRes[app_res.None] {
	if args == nil {
		return app_res.Reject[app_res.None](app_res.BadRequest, "分页参数不能为空")
	}

	if re := app_util.ClampOffsetLimit(args.Offset, &args.Limit); re.IsReject() {
		return re
	}

	return app_res.Accept(&app_res.None{})
}

// `vfyMarkReadSysMailId` validates mark-read message id.
func vfyMarkReadSysMailId(id string) app_res.AppRes[app_res.None] {
	if id == "" {
		return app_res.Reject[app_res.None](app_res.BadRequest, "id 不能为空")
	}

	return app_res.Accept(&app_res.None{})
}

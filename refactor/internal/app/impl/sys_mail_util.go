package app_impl

import (
	"poprako-s/internal/app/res"
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
func vfyListSysMailArgs(args *val.ListSysMailArgs) (res.ErrCode, string, bool) {
	if args == nil {
		return res.BadRequest, "分页参数不能为空", true
	}

	if code, msg, reject := app_util.ClampOffsetLimit(args.Offset, &args.Limit); reject {
		return code, msg, true
	}

	return 0, "", false
}

// `vfyMarkReadSysMailId` validates mark-read message id.
func vfyMarkReadSysMailId(id string) (res.ErrCode, string, bool) {
	if id == "" {
		return res.BadRequest, "id 不能为空", true
	}

	return 0, "", false
}

package app_impl

import (
	"strings"

	app_res "poprako-s/internal/app/res"
	app_util "poprako-s/internal/app/util"
	"poprako-s/internal/app/val"
	oss_iface "poprako-s/internal/domain/ext/oss"
	"poprako-s/internal/domain/model/aggr"
)

// `asmPageVal` converts a `Page` aggregate to app-facing `PageVal`.
func asmPageVal(page *aggr.Page, signer oss_iface.Signer) (*val.PageVal, error) {
	imageUrl := ""
	if page.ImageUploaded && page.ImageKey != nil && *page.ImageKey != "" {
		url, err := signer.GenGetUrl(*page.ImageKey)
		if err != nil {
			return nil, err
		}

		imageUrl = url
	}

	return &val.PageVal{
		Id:                  page.Id,
		ChapterId:           page.ChapterId,
		Index:               page.Index,
		ImageUrl:            imageUrl,
		ImageUploaded:       page.ImageUploaded,
		TotalUnitCount:      page.TotalUnitCount,
		TranslatedUnitCount: page.TranslatedUnitCount,
		ProofreadUnitCount:  page.ProofreadUnitCount,
		CreatedAt:           page.CreatedAt.UnixMilli(),
		UpdatedAt:           page.UpdatedAt.UnixMilli(),
	}, nil
}

// `vfyResvChapterPagesArgs` validates reservation arguments.
func vfyResvChapterPagesArgs(args *val.ResvChapterPagesArgs) app_res.AppRes[app_res.None] {
	if args == nil {
		return app_res.Reject[app_res.None](app_res.BadRequest, "预留参数不能为空")
	}
	if args.ChapterId == "" || args.PageCount <= 0 {
		return app_res.Reject[app_res.None](app_res.BadRequest, "chapter_id 和 page_count 必须合法")
	}

	args.FileExt = strings.TrimSpace(args.FileExt)
	if args.FileExt == "" {
		return app_res.Reject[app_res.None](app_res.BadRequest, "file_extension 不能为空")
	}

	return app_res.Accept(&app_res.None{})
}

// `vfyListChapterPageArgs` validates page list arguments.
func vfyListChapterPageArgs(args *val.ListChapterPageArgs) app_res.AppRes[app_res.None] {
	if args == nil {
		return app_res.Reject[app_res.None](app_res.BadRequest, "分页参数不能为空")
	}
	if args.ChapterId == "" {
		return app_res.Reject[app_res.None](app_res.BadRequest, "chapter_id 不能为空")
	}

	return app_util.ClampOffsetLimit(args.Offset, &args.Limit)
}

// `vfyMarkPageImageUploadedArgs` validates upload confirmation arguments.
func vfyMarkPageImageUploadedArgs(args *val.MarkPageImageUploadedArgs) app_res.AppRes[app_res.None] {
	if args == nil || args.PageId == "" {
		return app_res.Reject[app_res.None](app_res.BadRequest, "page_id 不能为空")
	}

	return app_res.Accept(&app_res.None{})
}

// `vfyDeleteByChapterId` validates chapter page deletion arguments.
func vfyDeleteByChapterId(chapterId string) app_res.AppRes[app_res.None] {
	if chapterId == "" {
		return app_res.Reject[app_res.None](app_res.BadRequest, "chapter_id 不能为空")
	}

	return app_res.Accept(&app_res.None{})
}

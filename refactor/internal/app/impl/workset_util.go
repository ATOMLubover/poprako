package app_impl

import (
	"poprako-s/internal/app/val"
	"poprako-s/internal/domain/model/aggr"
)

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

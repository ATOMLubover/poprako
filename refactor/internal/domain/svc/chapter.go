package svc

import (
	"fmt"

	"poprako-s/internal/domain/model/aggr"
	"poprako-s/pkg/util"
)

// `ChapterSvc` provides stateless helpers for chapter aggregate.
type ChapterSvc struct{}

// `NewChapterSvc` returns a ready-to-use `ChapterSvc`.
func NewChapterSvc() ChapterSvc {
	return ChapterSvc{}
}

// `NewChapterCre` builds a chapter creation payload with generated id.
func (ChapterSvc) NewChapterCre(
	comicId string,
	index int,
	subtitle *string,
	creatorId string,
) *aggr.ChapterCre {
	return &aggr.ChapterCre{
		Id:        util.GenId("chapter"),
		ComicId:   comicId,
		Index:     index,
		Subtitle:  subtitle,
		CreatorId: creatorId,
	}
}

// `DefSubtitle` returns default subtitle for chapter index.
func (ChapterSvc) DefSubtitle(index int) string {
	return fmt.Sprintf("Ch.%d", index)
}

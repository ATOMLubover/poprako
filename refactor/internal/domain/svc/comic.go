package svc

import (
	"poprako-s/internal/domain/model/aggr"
	"poprako-s/pkg/util"
)

// `ComicSvc` provides stateless helpers for comic aggregate workflows
// It only builds domain payloads and does not access repositories directly
// This keeps business construction logic centralized and reusable
type ComicSvc struct{}

// `NewComicSvc` returns a ready-to-use `ComicSvc`
func NewComicSvc() ComicSvc {
	return ComicSvc{}
}

// `NewComicCre` builds a `ComicCre` aggregate with generated id
// Caller should pre-compute `index` in transaction for consistency
// `desc` keeps nullable semantics to support put-style updates later
func (ComicSvc) NewComicCre(
	worksetId string,
	index int,
	title string,
	author string,
	desc *string,
	creatorId string,
) *aggr.ComicCre {
	return &aggr.ComicCre{
		Id:        util.GenId("comic"),
		WorksetId: worksetId,
		Index:     index,
		Title:     title,
		Author:    author,
		Desc:      desc,
		CreatorId: creatorId,
	}
}

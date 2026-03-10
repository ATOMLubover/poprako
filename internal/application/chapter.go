package application

import (
	"labelplus-next-web-be/internal/domain/external"
	"labelplus-next-web-be/internal/util"
	"labelplus-next-web-be/internal/value"

	"go.uber.org/zap"
)

type ChapterApplication interface{}

type chapterApplication struct {
	ossClient external.OSSClient
}

func NewChapterApplication(
	ossClient external.OSSClient,
) ChapterApplication {
	if ossClient == nil {
		zap.L().Panic(
			"NewChapterApplication: 依赖项不能为空",
			zap.Bool("ossClient_nil", ossClient == nil),
		)
	}

	return &chapterApplication{
		ossClient: ossClient,
	}
}

func (ca *chapterApplication) CreateChapter(
	scope util.TraceScope,
	currentUserID string,
	args value.CreateChapterArgs,
) (value.CreateChapterResult, error) {
}

func (ca *chapterApplication) ListComicChapters(
	scope util.TraceScope,
	currentUserID string,
	args value.ListComicChapterArgs,
) error {
}

func (ca *chapterApplication) UpdateChapter(
	scope util.TraceScope,
	currentUserID string,
	args value.UpdateChapterArgs,
) error {
}

func (ca *chapterApplication) DeleteComicChapter(
	scope util.TraceScope,
	currentUserID string,
	chapterID string,
) error {
}

package application

import (
	"errors"

	"labelplus-next-web-be/internal/application/adapter"
	"labelplus-next-web-be/internal/domain/model"
	"labelplus-next-web-be/internal/domain/repository"
	repository_infra "labelplus-next-web-be/internal/infrastructure/repository"
	"labelplus-next-web-be/internal/infrastructure/repository/query_option"
	"labelplus-next-web-be/internal/util"
	"labelplus-next-web-be/internal/value"

	"go.uber.org/zap"
)

type ChapterApplication interface {
	CreateComicChapter(
		scope util.TraceScope,
		currentUserID string,
		args value.CreateChapterArgs,
	) (value.CreateChapterResult, error)
	ListComicChapters(
		scope util.TraceScope,
		currentUserID string,
		args value.ListComicChapterArgs,
	) ([]value.ChapterDetail, error)
	UpdateChapter(
		scope util.TraceScope,
		currentUserID string,
		args value.UpdateChapterArgs,
	) error
	DeleteComicChapter(
		scope util.TraceScope,
		currentUserID string,
		chapterID string,
	) error
}

type chapterApplication struct {
	memberRepository  repository.MemberRepository
	comicRepository   repository.ComicRepository
	chapterRepository repository.ChapterRepository
}

func NewChapterApplication(
	memberRepository repository.MemberRepository,
	comicRepository repository.ComicRepository,
	chapterRepository repository.ChapterRepository,
) ChapterApplication {
	if memberRepository == nil ||
		comicRepository == nil ||
		chapterRepository == nil {
		zap.L().Panic(
			"NewChapterApplication: 依赖项不能为空",
			zap.Bool("memberRepository_nil", memberRepository == nil),
			zap.Bool("comicRepository_nil", comicRepository == nil),
			zap.Bool("chapterRepository_nil", chapterRepository == nil),
		)
	}

	return &chapterApplication{
		memberRepository:  memberRepository,
		comicRepository:   comicRepository,
		chapterRepository: chapterRepository,
	}
}

func (ca *chapterApplication) CreateComicChapter(
	scope util.TraceScope,
	currentUserID string,
	args value.CreateChapterArgs,
) (value.CreateChapterResult, error) {
	const fn = "ChapterApplication.CreateComicChapter"

	if err := args.Validate(); err != nil {
		scope.Logger().Warn(fn+": 参数验证失败", zap.Error(err))
		return value.CreateChapterResult{}, errors.New("参数错误: " + err.Error())
	}

	scope.
		WithFields(
			zap.String("current_user_id", currentUserID),
			zap.Any("args", args),
		).
		Logger().
		Debug(fn + ": 被调用")

	if !model.PermChapterCreate().Check(
		currentUserID,
		args.ComicID,
		adapter.HandleLoadMemberInfo(ca.memberRepository),
		adapter.HandleLoadComicInfo(ca.comicRepository),
	) {
		scope.Logger().Warn(fn + ": 权限检查失败")
		return value.CreateChapterResult{}, errors.New("权限不足")
	}

	transactionExecutor := ca.chapterRepository.BeginTransaction()
	if transactionExecutor.Error != nil {
		scope.Logger().Error(fn+": 开启事务失败", zap.Error(transactionExecutor.Error))
		return value.CreateChapterResult{}, errors.New("创建章节失败")
	}

	var transactionErr error

	defer func() {
		if transactionErr != nil {
			if rollbackErr := transactionExecutor.Rollback().Error; rollbackErr != nil {
				scope.Logger().Error(fn+": 回滚事务失败", zap.Error(rollbackErr))
			}
		}
	}()

	transactionErr = ca.chapterRepository.LockByComicID(transactionExecutor, args.ComicID)
	if transactionErr != nil {
		scope.Logger().Error(fn+": 锁定章节记录失败", zap.Error(transactionErr))
		return value.CreateChapterResult{}, errors.New("创建章节失败")
	}

	chapterCount, transactionErr := ca.chapterRepository.Count(
		transactionExecutor,
		query_option.ChapterQuery().FilterByComicID(args.ComicID),
	)
	if transactionErr != nil {
		scope.Logger().Error(fn+": 统计章节数量失败", zap.Error(transactionErr))
		return value.CreateChapterResult{}, errors.New("创建章节失败")
	}

	chapterCreation := model.NewChapterCreation(
		args.ComicID,
		int(chapterCount)+1,
		args.ChapterNo,
		currentUserID,
	)

	chapterID, transactionErr := ca.chapterRepository.Create(transactionExecutor, chapterCreation)
	if transactionErr != nil {
		scope.Logger().Error(fn+": 创建章节失败", zap.Error(transactionErr))
		return value.CreateChapterResult{}, errors.New("创建章节失败")
	}

	if commitErr := transactionExecutor.Commit().Error; commitErr != nil {
		scope.Logger().Error(fn+": 提交事务失败", zap.Error(commitErr))
		return value.CreateChapterResult{}, errors.New("创建章节失败")
	}

	return value.NewCreateChapterResultFromModel(chapterID), nil
}

func (ca *chapterApplication) ListComicChapters(
	scope util.TraceScope,
	currentUserID string,
	args value.ListComicChapterArgs,
) ([]value.ChapterDetail, error) {
	const fn = "ChapterApplication.ListComicChapters"

	if err := args.Validate(); err != nil {
		scope.Logger().Warn(fn+": 参数验证失败", zap.Error(err))
		return nil, errors.New("参数错误: " + err.Error())
	}

	scope.
		WithFields(
			zap.String("current_user_id", currentUserID),
			zap.Any("args", args),
		).
		Logger().
		Debug(fn + ": 被调用")

	if !model.PermChapterList().Check(
		currentUserID,
		args.ComicID,
		adapter.HandleLoadMemberInfo(ca.memberRepository),
		adapter.HandleLoadComicInfo(ca.comicRepository),
	) {
		scope.Logger().Warn(fn + ": 权限检查失败")
		return nil, errors.New("权限不足")
	}

	chapters, err := ca.chapterRepository.List(
		nil,
		query_option.ChapterQuery().FilterByComicID(args.ComicID),
		query_option.ChapterQuery().OrderByIndexAsc(),
		query_option.Paginate(args.Offset, args.Limit),
	)
	if err != nil {
		scope.Logger().Error(fn+": 获取章节列表失败", zap.Error(err))
		return nil, errors.New("获取章节列表失败")
	}

	result := make([]value.ChapterDetail, len(chapters))
	for i, chapter := range chapters {
		result[i] = value.NewChapterInfoFromModel(chapter)
	}

	return result, nil
}

func (ca *chapterApplication) UpdateChapter(
	scope util.TraceScope,
	currentUserID string,
	args value.UpdateChapterArgs,
) error {
	const fn = "ChapterApplication.UpdateChapter"

	if err := args.Validate(); err != nil {
		scope.Logger().Warn(fn+": 参数验证失败", zap.Error(err))
		return errors.New("参数错误: " + err.Error())
	}

	scope.
		WithFields(
			zap.String("current_user_id", currentUserID),
			zap.Any("args", args),
		).
		Logger().
		Debug(fn + ": 被调用")

	targetChapter, err := ca.chapterRepository.Get(
		nil,
		query_option.FilterByID(repository_infra.ChapterTable, args.ChapterID),
	)
	if err != nil {
		scope.Logger().Error(fn+": 获取目标章节信息失败", zap.Error(err))
		return errors.New("无法获取章节信息")
	}

	if !model.PermChapterUpdate().Check(
		currentUserID,
		targetChapter.ComicID,
		adapter.HandleLoadMemberInfo(ca.memberRepository),
		adapter.HandleLoadComicInfo(ca.comicRepository),
	) {
		scope.Logger().Warn(fn + ": 权限检查失败")
		return errors.New("权限不足")
	}

	var chapterNo *string
	if args.ChapterNo != "" {
		chapterNo = &args.ChapterNo
	}

	chapterUpdate := model.NewChapterUpdate(args.ChapterID, chapterNo)

	if err := ca.chapterRepository.Update(nil, chapterUpdate); err != nil {
		scope.Logger().Error(fn+": 更新章节失败", zap.Error(err))
		return errors.New("更新章节失败")
	}

	return nil
}

func (ca *chapterApplication) DeleteComicChapter(
	scope util.TraceScope,
	currentUserID string,
	chapterID string,
) error {
	const fn = "ChapterApplication.DeleteComicChapter"

	if chapterID == "" {
		scope.Logger().Warn(fn + ": chapterID 为空")
		return errors.New(ErrInternalError)
	}

	scope.
		WithFields(
			zap.String("current_user_id", currentUserID),
			zap.String("chapter_id", chapterID),
		).
		Logger().
		Debug(fn + ": 被调用")

	targetChapter, err := ca.chapterRepository.Get(
		nil,
		query_option.FilterByID(repository_infra.ChapterTable, chapterID),
	)
	if err != nil {
		scope.Logger().Error(fn+": 获取目标章节信息失败", zap.Error(err))
		return errors.New("无法获取章节信息")
	}

	if !model.PermChapterDelete().Check(
		currentUserID,
		targetChapter.ComicID,
		adapter.HandleLoadMemberInfo(ca.memberRepository),
		adapter.HandleLoadComicInfo(ca.comicRepository),
	) {
		scope.Logger().Warn(fn + ": 权限检查失败")
		return errors.New("权限不足")
	}

	if err := ca.chapterRepository.Delete(nil, chapterID); err != nil {
		scope.Logger().Error(fn+": 删除章节失败", zap.Error(err))
		return errors.New("删除章节失败")
	}

	return nil
}

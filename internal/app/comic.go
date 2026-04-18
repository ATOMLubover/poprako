package app

import (
	"context"
	"errors"
	"path/filepath"

	event_handler "poprako-s/internal/app/event_handler"
	"poprako-s/internal/app/val"
	"poprako-s/internal/domain/event"
	"poprako-s/internal/domain/ext/oss"
	"poprako-s/internal/domain/model"
	"poprako-s/internal/domain/repo"
	"poprako-s/internal/domain/service"

	"go.uber.org/zap"
)

type ComicApp interface {
	// List 获取指定作品集的漫画列表
	List(
		cx context.Context,
		currUserID string,
		args *val.ListComicArgs,
	) ([]*val.ComicInfo, error)

	// Create 创建一部新漫画
	Create(
		cx context.Context,
		currUserID string,
		args *val.CreateComicArgs,
	) (*val.CreateComicRes, error)

	// Update 更新漫画信息
	Update(
		cx context.Context,
		currUserID string,
		args *val.UpdateComicArgs,
	) error

	// Remove 删除漫画
	Remove(
		cx context.Context,
		currUserID string,
		comicID string,
	) error

	// ReserveCover 为漫画封面生成预签名上传 URL，并预留 cover_oss_key
	ReserveCover(
		cx context.Context,
		currUserID string,
		args *val.ReserveComicCoverArgs,
	) (*val.ReserveComicCoverRes, error)

	// ConfirmCoverUploaded 确认漫画封面已完成上传
	ConfirmCoverUploaded(
		cx context.Context,
		currUserID string,
		comicID string,
	) error
}

type comicAppImpl struct {
	comicSvc service.ComicService

	memberRepo  repo.MemberRepo
	worksetRepo repo.WorksetRepo
	comicRepo   repo.ComicRepo
	txnMgr      repo.TxnMgr
	eventBus    event.EventBus
	ossClient   oss.Client
}

func NewComicApp(
	comicSvc service.ComicService,
	memberRepo repo.MemberRepo,
	worksetRepo repo.WorksetRepo,
	comicRepo repo.ComicRepo,
	txnMgr repo.TxnMgr,
	eventBus event.EventBus,
	ossClient oss.Client,
) ComicApp {
	// 校验构造函数依赖
	if comicSvc == nil ||
		memberRepo == nil ||
		worksetRepo == nil ||
		comicRepo == nil ||
		txnMgr == nil ||
		eventBus == nil ||
		ossClient == nil {
		zap.L().Panic(
			"NewComicApp: 依赖项不能为空",
			zap.Bool("comicSvc_nil", comicSvc == nil),
			zap.Bool("memberRepo_nil", memberRepo == nil),
			zap.Bool("worksetRepo_nil", worksetRepo == nil),
			zap.Bool("comicRepo_nil", comicRepo == nil),
			zap.Bool("txnMgr_nil", txnMgr == nil),
			zap.Bool("eventBus_nil", eventBus == nil),
			zap.Bool("ossClient_nil", ossClient == nil),
		)
	}

	// 返回真实业务实现
	return &comicAppImpl{
		comicSvc:    comicSvc,
		memberRepo:  memberRepo,
		worksetRepo: worksetRepo,
		comicRepo:   comicRepo,
		txnMgr:      txnMgr,
		eventBus:    eventBus,
		ossClient:   ossClient,
	}
}

func (a *comicAppImpl) List(
	cx context.Context,
	currUserID string,
	args *val.ListComicArgs,
) ([]*val.ComicInfo, error) {
	// 获取上下文中的日志记录器
	lgr := retrieveLgr(cx)

	// 通过作品集获取所属汉化组 ID 用于鉴权
	targetWorkset, err := a.worksetRepo.GetByID(args.WorksetID)
	if err != nil {
		// 记录查询失败
		lgr.Error(
			"获取漫画列表失败：获取作品集信息失败",
			zap.String("workset_id", args.WorksetID),
			zap.Error(err),
		)

		// 返回客户端可展示的错误
		return nil, errors.New("无法获取作品集信息")
	}

	// 鉴权：检查当前用户是否为该汉化组成员
	_, err = a.memberRepo.Get(model.MemberQueryOpt{
		UserID: &currUserID,
		TeamID: &targetWorkset.TeamID,
	})
	if err != nil {
		// 记录权限校验失败
		lgr.Warn(
			"获取漫画列表失败：权限不足",
			zap.String("curr_user_id", currUserID),
		)

		// 返回客户端可展示的错误
		return nil, errors.New("权限不足")
	}

	// 构造查询条件
	queryOpt := model.ComicQueryOpt{
		WorksetID:       args.WorksetID,
		UploadStatus:    args.UploadStatus,
		TranslateStatus: args.TranslateStatus,
		ProofreadStatus: args.ProofreadStatus,
		TypesetStatus:   args.TypesetStatus,
		ReviewStatus:    args.ReviewStatus,
		PublishStatus:   args.PublishStatus,
		Includes:        args.Includes,
	}

	if args.FuzzyTitle != "" {
		queryOpt.FuzzyTitle = &args.FuzzyTitle
	}

	// 查询漫画列表
	comics, err := a.comicRepo.List(queryOpt)
	if err != nil {
		// 记录查询失败
		lgr.Error(
			"获取漫画列表失败",
			zap.String("workset_id", args.WorksetID),
			zap.Error(err),
		)

		// 返回客户端可展示的错误
		return nil, errors.New("获取漫画列表失败")
	}

	// 组装为 app 层值对象列表
	result := make([]*val.ComicInfo, len(comics))

	for i, comic := range comics {
		result[i] = assembleComicInfo(&comic, a.ossClient)
	}

	// 返回漫画列表
	return result, nil
}

func (a *comicAppImpl) Create(
	cx context.Context,
	currUserID string,
	args *val.CreateComicArgs,
) (*val.CreateComicRes, error) {
	// 获取上下文中的日志记录器
	lgr := retrieveLgr(cx)

	// 在事务中创建漫画（需要 count 获取 index）

	var createdID string

	if err := a.txnMgr.RunInTxn(func(cx context.Context) error {
		comicRepoTxn, err := a.comicRepo.FromTxnCx(cx)
		if err != nil {
			return err
		}

		memberRepoTxn, err := a.memberRepo.FromTxnCx(cx)
		if err != nil {
			return err
		}

		worksetRepoTxn, err := a.worksetRepo.FromTxnCx(cx)
		if err != nil {
			return err
		}

		// 统计当前作品集下的漫画数量以确定 index
		count, err := comicRepoTxn.Count(model.ComicQueryOpt{
			WorksetID: args.WorksetID,
		})
		if err != nil {
			return err
		}

		// 通过领域服务构造创建载荷（含权限校验）
		creation, err := a.comicSvc.NewCreation(
			memberRepoTxn,
			worksetRepoTxn,
			currUserID,
			args.WorksetID,
			int(count),
			args.Title,
			args.Author,
			args.Description,
			currUserID,
		)
		if err != nil {
			return err
		}

		// 持久化漫画
		comicInfo, err := comicRepoTxn.Create(creation)
		if err != nil {
			return err
		}

		createdID = comicInfo.ID
		eventCx := event_handler.WithWorksetRepoTxn(cx, worksetRepoTxn)

		return a.eventBus.Pub(eventCx, creation.PullEvents())
	}); err != nil {
		// 记录创建失败
		lgr.Error(
			"创建漫画失败",
			zap.String("curr_user_id", currUserID),
			zap.String("workset_id", args.WorksetID),
			zap.Error(err),
		)

		// 返回客户端可展示的错误
		return nil, errors.New("创建漫画失败")
	}

	// 返回创建结果
	return &val.CreateComicRes{ID: createdID}, nil
}

func (a *comicAppImpl) Update(
	cx context.Context,
	currUserID string,
	args *val.UpdateComicArgs,
) error {
	// 获取上下文中的日志记录器
	lgr := retrieveLgr(cx)

	// 查询目标漫画信息以获取所属作品集
	targetComic, err := a.comicRepo.GetByID(args.ID)
	if err != nil {
		// 记录查询失败
		lgr.Error(
			"更新漫画失败：获取目标漫画信息失败",
			zap.String("comic_id", args.ID),
			zap.Error(err),
		)

		// 返回客户端可展示的错误
		return errors.New("无法获取漫画信息")
	}

	// 通过作品集获取所属汉化组 ID 用于鉴权
	targetWorkset, err := a.worksetRepo.GetByID(targetComic.WorksetID)
	if err != nil {
		// 记录查询失败
		lgr.Error(
			"更新漫画失败：获取作品集信息失败",
			zap.String("workset_id", targetComic.WorksetID),
			zap.Error(err),
		)

		// 返回客户端可展示的错误
		return errors.New("无法获取作品集信息")
	}

	// 鉴权：检查当前用户在漫画所属汉化组中是否为管理员
	currMember, err := a.memberRepo.Get(model.MemberQueryOpt{
		UserID: &currUserID,
		TeamID: &targetWorkset.TeamID,
	})
	if err != nil || !currMember.HasAnyRole(model.RoleAdmin) {
		// 记录权限校验失败
		lgr.Warn(
			"更新漫画失败：权限不足",
			zap.String("curr_user_id", currUserID),
		)

		// 返回客户端可展示的错误
		return errors.New("权限不足")
	}

	// 构造更新载荷
	update := &model.ComicUpdate{
		ID:          args.ID,
		Title:       args.Title,
		Author:      args.Author,
		Description: args.Description,
	}

	// 持久化更新
	if err := a.comicRepo.Update(update); err != nil {
		// 记录更新失败
		lgr.Error(
			"更新漫画失败",
			zap.String("comic_id", args.ID),
			zap.Error(err),
		)

		// 返回客户端可展示的错误
		return errors.New("更新漫画失败")
	}

	// 返回更新成功
	return nil
}

func (a *comicAppImpl) Remove(
	cx context.Context,
	currUserID string,
	comicID string,
) error {
	// 获取上下文中的日志记录器
	lgr := retrieveLgr(cx)

	// 查询目标漫画信息以获取所属作品集
	targetComic, err := a.comicRepo.GetByID(comicID)
	if err != nil {
		// 记录查询失败
		lgr.Error(
			"删除漫画失败：获取目标漫画信息失败",
			zap.String("comic_id", comicID),
			zap.Error(err),
		)

		// 返回客户端可展示的错误
		return errors.New("无法获取漫画信息")
	}

	// 通过作品集获取所属汉化组 ID 用于鉴权
	targetWorkset, err := a.worksetRepo.GetByID(targetComic.WorksetID)
	if err != nil {
		// 记录查询失败
		lgr.Error(
			"删除漫画失败：获取作品集信息失败",
			zap.String("workset_id", targetComic.WorksetID),
			zap.Error(err),
		)

		// 返回客户端可展示的错误
		return errors.New("无法获取作品集信息")
	}

	// 鉴权：检查当前用户在漫画所属汉化组中是否为管理员
	currMember, err := a.memberRepo.Get(model.MemberQueryOpt{
		UserID: &currUserID,
		TeamID: &targetWorkset.TeamID,
	})
	if err != nil || !currMember.HasAnyRole(model.RoleAdmin) {
		// 记录权限校验失败
		lgr.Warn(
			"删除漫画失败：权限不足",
			zap.String("curr_user_id", currUserID),
		)

		// 返回客户端可展示的错误
		return errors.New("权限不足")
	}

	if err := newOSSDeleteExecutor(a.ossClient).deleteOne(targetComic.CoverOSSKey); err != nil {
		lgr.Error(
			"删除漫画失败：删除封面 OSS 资源失败",
			zap.String("comic_id", comicID),
			zap.String("cover_oss_key", targetComic.CoverOSSKey),
			zap.Error(err),
		)

		return errors.New("删除漫画失败")
	}

	if err := a.txnMgr.RunInTxn(func(cx context.Context) error {
		comicRepoTxn, err := a.comicRepo.FromTxnCx(cx)
		if err != nil {
			return err
		}

		worksetRepoTxn, err := a.worksetRepo.FromTxnCx(cx)
		if err != nil {
			return err
		}

		if err := comicRepoTxn.Delete(comicID); err != nil {
			return err
		}

		eventCx := event_handler.WithWorksetRepoTxn(cx, worksetRepoTxn)

		removalEvent := a.comicSvc.NewRemovalEvent(comicID, targetComic.WorksetID)

		return a.eventBus.Pub(eventCx, []event.Event{removalEvent})
	}); err != nil {
		lgr.Error(
			"删除漫画失败",
			zap.String("comic_id", comicID),
			zap.Error(err),
		)

		return errors.New("删除漫画失败")
	}

	return nil
}

func (a *comicAppImpl) ReserveCover(
	cx context.Context,
	currUserID string,
	args *val.ReserveComicCoverArgs,
) (*val.ReserveComicCoverRes, error) {
	// 获取上下文中的日志记录器
	lgr := retrieveLgr(cx)

	// 获取漫画信息以解析所属汉化组，用于鉴权
	targetComic, err := a.comicRepo.GetByID(args.ComicID)
	if err != nil {
		// 记录查询失败
		lgr.Warn(
			"预留漫画封面失败：查询漫画信息失败",
			zap.String("comic_id", args.ComicID),
			zap.Error(err),
		)

		// 返回客户端可展示的错误
		return nil, errors.New("漫画不存在")
	}

	// 通过作品集获取所属汉化组 ID 用于鉴权
	targetWorkset, err := a.worksetRepo.GetByID(targetComic.WorksetID)
	if err != nil {
		// 记录查询失败
		lgr.Warn(
			"预留漫画封面失败：查询作品集信息失败",
			zap.String("workset_id", targetComic.WorksetID),
			zap.Error(err),
		)

		// 返回客户端可展示的错误
		return nil, errors.New("无法获取作品集信息")
	}

	// 鉴权：检查当前用户是否为汉化组管理员
	currMember, err := a.memberRepo.Get(model.MemberQueryOpt{
		UserID: &currUserID,
		TeamID: &targetWorkset.TeamID,
	})
	if err != nil || !currMember.HasAnyRole(model.RoleAdmin) {
		// 记录权限校验失败
		lgr.Warn(
			"预留漫画封面失败：权限不足",
			zap.String("curr_user_id", currUserID),
			zap.String("comic_id", args.ComicID),
		)

		// 返回客户端可展示的错误
		return nil, errors.New("权限不足")
	}

	// 提取文件扩展名，与 OSS Key 拼接以便 OSS 正确识别 Content-Type
	ext := filepath.Ext(args.FileName)

	// 基于漫画 ID 生成封面对象 Key，并附加扩展名
	coverOSSKey := a.comicSvc.GenCoverOSSKey(args.ComicID) + ext

	// 为客户端生成预签名上传链接
	putURL, err := a.ossClient.GeneratePutPresignedURL(coverOSSKey)
	if err != nil {
		// 记录上传链接生成失败
		lgr.Error(
			"预留漫画封面失败：生成上传链接失败",
			zap.String("comic_id", args.ComicID),
			zap.Error(err),
		)

		// 返回客户端可展示的错误
		return nil, errors.New("预留漫画封面失败")
	}

	// 在数据库中预填充封面对象 Key
	if err := a.comicRepo.PreFillCoverOSSKey(args.ComicID, coverOSSKey); err != nil {
		// 记录预写失败
		lgr.Error(
			"预留漫画封面失败：写入封面 OSS Key 失败",
			zap.String("comic_id", args.ComicID),
			zap.Error(err),
		)

		// 返回客户端可展示的错误
		return nil, errors.New("预留漫画封面失败")
	}

	// 返回预留结果
	return &val.ReserveComicCoverRes{PutURL: putURL}, nil
}

func (a *comicAppImpl) ConfirmCoverUploaded(
	cx context.Context,
	currUserID string,
	comicID string,
) error {
	// 获取上下文中的日志记录器
	lgr := retrieveLgr(cx)

	// 获取漫画信息以解析所属汉化组，用于鉴权
	targetComic, err := a.comicRepo.GetByID(comicID)
	if err != nil {
		// 记录查询失败
		lgr.Warn(
			"确认漫画封面上传失败：查询漫画信息失败",
			zap.String("comic_id", comicID),
			zap.Error(err),
		)

		// 返回客户端可展示的错误
		return errors.New("漫画不存在")
	}

	// 通过作品集获取所属汉化组 ID 用于鉴权
	targetWorkset, err := a.worksetRepo.GetByID(targetComic.WorksetID)
	if err != nil {
		// 记录查询失败
		lgr.Warn(
			"确认漫画封面上传失败：查询作品集信息失败",
			zap.String("workset_id", targetComic.WorksetID),
			zap.Error(err),
		)

		// 返回客户端可展示的错误
		return errors.New("无法获取作品集信息")
	}

	// 鉴权：检查当前用户是否为汉化组管理员
	currMember, err := a.memberRepo.Get(model.MemberQueryOpt{
		UserID: &currUserID,
		TeamID: &targetWorkset.TeamID,
	})
	if err != nil || !currMember.HasAnyRole(model.RoleAdmin) {
		// 记录权限校验失败
		lgr.Warn(
			"确认漫画封面上传失败：权限不足",
			zap.String("curr_user_id", currUserID),
			zap.String("comic_id", comicID),
		)

		// 返回客户端可展示的错误
		return errors.New("权限不足")
	}

	// 将封面状态标记为已上传
	if err := a.comicRepo.ConfirmCoverUploaded(comicID); err != nil {
		// 记录确认失败
		lgr.Error(
			"确认漫画封面上传失败",
			zap.String("comic_id", comicID),
			zap.Error(err),
		)

		// 返回客户端可展示的错误
		return errors.New("确认漫画封面上传失败")
	}

	// 返回确认成功
	return nil
}

// assembleComicInfo 将领域层漫画信息转换为 app 层值对象
func assembleComicInfo(
	info *model.ComicInfo,
	ossClient oss.Client,
) *val.ComicInfo {
	result := &val.ComicInfo{
		ID:              info.ID,
		WorksetID:       info.WorksetID,
		Index:           info.Index,
		Title:           info.Title,
		Author:          info.Author,
		Description:     info.Description,
		ChapterCount:    info.ChapterCount,
		IsCoverUploaded: info.IsCoverUploaded,
		CreatorID:       info.CreatorID,
		LastActiveAt:    info.LastActiveAt.UnixMilli(),
		CreatedAt:       info.CreatedAt.UnixMilli(),
		UpdatedAt:       info.UpdatedAt.UnixMilli(),
	}

	// 若封面已上传则生成可访问地址
	if info.IsCoverUploaded && info.CoverOSSKey != "" {
		if coverURL, err := ossClient.GenerateGetPresignedURL(info.CoverOSSKey); err == nil {
			result.CoverURL = coverURL
		}
	}

	// 若包含创建者信息则一并组装
	if info.Creator != nil {
		userInfo, _ := assembleUserInfo(info.Creator, ossClient)
		result.Creator = userInfo
	}

	// 若包含作品集信息则一并组装
	if info.Workset != nil {
		result.Workset = assembleWorksetInfo(info.Workset)
	}

	return result
}

// logComicAppImpl 是 ComicApp 的日志包装实现
type logComicAppImpl struct {
	app ComicApp
}

func NewLogComicApp(
	app ComicApp,
) ComicApp {
	if app == nil {
		zap.L().Panic(
			"NewLogComicApp: 依赖项不能为空",
			zap.Bool("app_nil", app == nil),
		)
	}

	return &logComicAppImpl{app: app}
}

func (a *logComicAppImpl) List(
	cx context.Context,
	currUserID string,
	args *val.ListComicArgs,
) ([]*val.ComicInfo, error) {
	if a == nil || a.app == nil {
		return nil, errors.New("ComicApp 不可用")
	}

	if args == nil || args.WorksetID == "" {
		return nil, errors.New("参数不合法")
	}

	lgr := retrieveLgr(cx).With(zap.String("method", "List"), zap.String("curr_user_id", currUserID), zap.String("workset_id", args.WorksetID))

	cx = injectLgr(cx, lgr)

	lgr.Info("[logComicAppImpl.List] CALL")

	return a.app.List(cx, currUserID, args)
}

func (a *logComicAppImpl) Create(
	cx context.Context,
	currUserID string,
	args *val.CreateComicArgs,
) (*val.CreateComicRes, error) {
	if a == nil || a.app == nil {
		return nil, errors.New("ComicApp 不可用")
	}

	if args == nil || args.WorksetID == "" || args.Title == "" {
		return nil, errors.New("参数不合法")
	}

	lgr := retrieveLgr(cx).With(zap.String("method", "Create"), zap.String("curr_user_id", currUserID))

	cx = injectLgr(cx, lgr)

	lgr.Info("[logComicAppImpl.Create] CALL")

	return a.app.Create(cx, currUserID, args)
}

func (a *logComicAppImpl) Update(
	cx context.Context,
	currUserID string,
	args *val.UpdateComicArgs,
) error {
	if a == nil || a.app == nil {
		return errors.New("ComicApp 不可用")
	}

	if args == nil || args.ID == "" {
		return errors.New("参数不合法")
	}

	lgr := retrieveLgr(cx).With(zap.String("method", "Update"), zap.String("curr_user_id", currUserID), zap.String("comic_id", args.ID))

	cx = injectLgr(cx, lgr)

	lgr.Info("[logComicAppImpl.Update] CALL")

	return a.app.Update(cx, currUserID, args)
}

func (a *logComicAppImpl) Remove(
	cx context.Context,
	currUserID string,
	comicID string,
) error {
	if a == nil || a.app == nil {
		return errors.New("ComicApp 不可用")
	}

	if comicID == "" {
		return errors.New("漫画 ID 不能为空")
	}

	lgr := retrieveLgr(cx).With(zap.String("method", "Remove"), zap.String("curr_user_id", currUserID), zap.String("comic_id", comicID))

	cx = injectLgr(cx, lgr)

	lgr.Info("[logComicAppImpl.Remove] CALL")

	return a.app.Remove(cx, currUserID, comicID)
}

func (a *logComicAppImpl) ReserveCover(
	cx context.Context,
	currUserID string,
	args *val.ReserveComicCoverArgs,
) (*val.ReserveComicCoverRes, error) {
	if a == nil || a.app == nil {
		return nil, errors.New("ComicApp 不可用")
	}

	if args == nil || args.ComicID == "" {
		return nil, errors.New("漫画 ID 不能为空")
	}

	if args.FileName == "" {
		return nil, errors.New("文件名不能为空")
	}

	lgr := retrieveLgr(cx).With(zap.String("method", "ReserveCover"), zap.String("curr_user_id", currUserID), zap.String("comic_id", args.ComicID))

	cx = injectLgr(cx, lgr)

	lgr.Info("[logComicAppImpl.ReserveCover] CALL")

	return a.app.ReserveCover(cx, currUserID, args)
}

func (a *logComicAppImpl) ConfirmCoverUploaded(
	cx context.Context,
	currUserID string,
	comicID string,
) error {
	if a == nil || a.app == nil {
		return errors.New("ComicApp 不可用")
	}

	if comicID == "" {
		return errors.New("漫画 ID 不能为空")
	}

	lgr := retrieveLgr(cx).With(zap.String("method", "ConfirmCoverUploaded"), zap.String("curr_user_id", currUserID), zap.String("comic_id", comicID))

	cx = injectLgr(cx, lgr)

	lgr.Info("[logComicAppImpl.ConfirmCoverUploaded] CALL")

	return a.app.ConfirmCoverUploaded(cx, currUserID, comicID)
}

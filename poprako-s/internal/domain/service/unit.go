package service

import (
	"errors"

	"poprako-s/internal/domain/model"
	"poprako-s/internal/domain/repo"
)

type UnitService interface {
	// NewDiff 根据业务参数创建 UnitDiff 领域模型
	// 只有 currUser 是当前 page 的 chapter 的翻译、校对时才能创建 UnitDiff，否则返回错误
	// 另外和校对相关的字段，也只能由校对创建或修改，翻译无权操作
	NewDiff(
		pr repo.PageRepo,
		ar repo.AssignmentRepo,
		currUser *model.UserInfo,
		pageID string,
		insert []model.UnitCreation,
		patch []model.UnitPatch,
		delete []string,
	) (*model.UnitDiff, error)
}

// unitServiceImpl 是 UnitService 的具体实现 无内禀状态
type unitServiceImpl struct{}

func NewUnitService() UnitService {
	// 返回无状态实现
	return &unitServiceImpl{}
}

func (s *unitServiceImpl) NewDiff(
	pr repo.PageRepo,
	ar repo.AssignmentRepo,
	currUser *model.UserInfo,
	pageID string,
	insert []model.UnitCreation,
	patch []model.UnitPatch,
	delete []string,
) (*model.UnitDiff, error) {
	// 查询页面信息
	page, err := pr.GetByID(pageID)
	if err != nil {
		// 返回页面不存在错误
		return nil, errors.New("页面不存在，无法提交 unit 变更")
	}

	// 查询当前用户在章节中的分配信息
	assignment, err := ar.Get(model.AssignmentQueryOpt{
		ChapterID: &page.ChapterID,
		UserID:    &currUser.ID,
	})
	if err != nil {
		// 返回权限错误
		return nil, errors.New("仅章节翻译或校对可以提交 unit 变更")
	}

	// 计算当前用户角色
	isTranslator := assignment.HasAnyRole(model.RoleTranslator)

	isProofreader := assignment.HasAnyRole(model.RoleProofreader)

	if !isTranslator && !isProofreader {
		// 返回权限错误
		return nil, errors.New("仅章节翻译或校对可以提交 unit 变更")
	}

	// 翻译无权操作校对字段
	if isTranslator && !isProofreader {
		for _, u := range insert {
			if touchesProofreadFieldsInInsert(u) {
				// 返回字段权限错误
				return nil, errors.New("翻译无权创建或修改校对相关字段")
			}
		}

		for _, u := range patch {
			if touchesProofreadFieldsInPatch(u) {
				// 返回字段权限错误
				return nil, errors.New("翻译无权创建或修改校对相关字段")
			}
		}
	}

	// 返回差异结果
	return &model.UnitDiff{
		Insert: insert,
		Patch:  patch,
		Delete: delete,
	}, nil
}

func touchesProofreadFieldsInInsert(
	u model.UnitCreation,
) bool {
	// 依次检查校对字段
	if u.IsProofread {
		return true
	}

	if u.ProofreadText != nil {
		return true
	}

	if u.ProofreaderID != nil {
		return true
	}

	if u.ProofreaderComment != nil {
		return true
	}

	// 未触碰时返回 false
	return false
}

func touchesProofreadFieldsInPatch(
	u model.UnitPatch,
) bool {
	// 依次检查校对字段
	if u.IsProofread != nil {
		return true
	}

	if u.ProofreadText != nil {
		return true
	}

	if u.ProofreaderID != nil {
		return true
	}

	if u.ProofreaderComment != nil {
		return true
	}

	// 未触碰时返回 false
	return false
}

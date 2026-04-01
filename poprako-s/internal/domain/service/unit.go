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
		currUser *model.UserInfo,
		pageID string,
		insert []model.UnitCreation,
		patch []model.UnitPatch,
		delete []string,
	) (*model.UnitDiff, error)
}

type unitServiceImpl struct {
	pageRepo       repo.PageRepo
	assignmentRepo repo.AssignmentRepo
}

func NewUnitService(pageRepo repo.PageRepo, assignmentRepo repo.AssignmentRepo) UnitService {
	return &unitServiceImpl{
		pageRepo:       pageRepo,
		assignmentRepo: assignmentRepo,
	}
}

func (s *unitServiceImpl) NewDiff(
	currUser *model.UserInfo,
	pageID string,
	insert []model.UnitCreation,
	patch []model.UnitPatch,
	delete []string,
) (*model.UnitDiff, error) {
	page, err := s.pageRepo.GetByID(pageID)
	if err != nil {
		return nil, errors.New("页面不存在，无法提交 unit 变更")
	}

	assignment, err := s.assignmentRepo.Get(model.AssignmentQueryOpt{
		ChapterID: &page.ChapterID,
		UserID:    &currUser.ID,
	})
	if err != nil {
		return nil, errors.New("仅章节翻译或校对可以提交 unit 变更")
	}

	isTranslator := assignment.HasAnyRole(model.RoleTranslator)
	isProofreader := assignment.HasAnyRole(model.RoleProofreader)
	if !isTranslator && !isProofreader {
		return nil, errors.New("仅章节翻译或校对可以提交 unit 变更")
	}

	if isTranslator && !isProofreader {
		for _, u := range insert {
			if touchesProofreadFieldsInInsert(u) {
				return nil, errors.New("翻译无权创建或修改校对相关字段")
			}
		}
		for _, u := range patch {
			if touchesProofreadFieldsInPatch(u) {
				return nil, errors.New("翻译无权创建或修改校对相关字段")
			}
		}
	}

	return &model.UnitDiff{
		Insert: insert,
		Patch:  patch,
		Delete: delete,
	}, nil
}

func touchesProofreadFieldsInInsert(u model.UnitCreation) bool {
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

	return false
}

func touchesProofreadFieldsInPatch(u model.UnitPatch) bool {
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

	return false
}

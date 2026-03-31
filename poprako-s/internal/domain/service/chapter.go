package service

import (
	"fmt"

	"poprako-s/internal/domain/model"
)

type ChapterService interface {
	// TransiteWorkflow 接受一个工作流转换事件，根据事件类型和当前状态执行相应的状态转换
	// 它负责权限检验，只有当用户 u 有权执行事件 t 时才会执行状态转换，否则返回错误
	TransiteWorkflow(t model.WorkflowTransition, c *model.ChapterInfo, a *model.AssignmentInfo) error
}

type chapterServiceImpl struct{}

func NewChapterService() ChapterService {
	return &chapterServiceImpl{}
}

// TransiteWorkflow 接受一个工作流转换事件，根据事件类型和当前状态执行相应的状态转换
// 它负责权限检验，只有当用户 u 有权执行事件 t 时才会执行状态转换，否则返回错误
func (s *chapterServiceImpl) TransiteWorkflow(
	t model.WorkflowTransition,
	c *model.ChapterInfo,
	a *model.AssignmentInfo,
) error {
	if a.HasAnyRole(model.RoleReviewer) {
		// 监修可以执行任何转换
		return c.TransiteWorkflow(t)
	}

	switch t {
	case model.WorkflowUploadComplete:
		if !a.HasAnyRole(model.RoleRawProvider) {
			return fmt.Errorf("只有图源可以标记上传完成")
		}
		return c.TransiteWorkflow(t)

	case model.WorkflowTranslateStart, model.WorkflowTranslateComplete:
		if !a.HasAnyRole(model.RoleTranslator) {
			return fmt.Errorf("只有翻译可以标记翻译开始或完成")
		}
		return c.TransiteWorkflow(t)

	case model.WorkflowProofreadStart, model.WorkflowProofreadComplete:
		if !a.HasAnyRole(model.RoleProofreader) {
			return fmt.Errorf("只有校对可以标记校对开始或完成")
		}
		return c.TransiteWorkflow(t)

	case model.WorkflowTypesetStart, model.WorkflowTypesetComplete:
		if !a.HasAnyRole(model.RoleTypesetter) {
			return fmt.Errorf("只有嵌字可以标记嵌字开始或完成")
		}
		return c.TransiteWorkflow(t)

	case model.WorkflowReviewComplete:
		if !a.HasAnyRole(model.RoleReviewer) {
			return fmt.Errorf("只有监修可以标记监修完成")
		}
		return c.TransiteWorkflow(t)

	case model.WorkflowPublishComplete:
		if !a.HasAnyRole(model.RolePublisher) {
			return fmt.Errorf("只有发布可以标记发布完成")
		}
		return c.TransiteWorkflow(t)

	default:
		return fmt.Errorf("未知的工作流转换：%s", t)
	}
}

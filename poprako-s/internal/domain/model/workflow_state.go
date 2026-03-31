package model

import (
	"fmt"
	"time"
)

// workflowState 内部维护所有进度相关信息
type workflowState struct {
	UploadStatus    UploadStatus
	TranslateStatus TranslateStatus
	ProofreadStatus ProofreadStatus
	TypesettStatus  TypesettStatus
	ReviewStatus    ReviewStatus
	PublishStatus   PublishStatus
}

// Transite 接受一个工作流转换事件，根据事件类型和当前状态执行相应的状态转换
// 它不负责权限检验，调用者必须保证事件的合法性
func (s *workflowState) transite(t WorkflowTransition) error {
	now := time.Now()

	switch t {
	// ========================
	// Upload
	// ========================
	case WorkflowUploadComplete:
		switch s.UploadStatus.Status {
		case WorkflowOngoing:
			s.UploadStatus.Status = WorkflowCompleted
			s.UploadStatus.CompletedAt = &now
		case WorkflowPending:
			return fmt.Errorf("上传进度尚未标记为进行中，无法标记为完成")
		case WorkflowCompleted:
			return fmt.Errorf("上传进度已标记为已完成，无法重复标记")
		default:
			return fmt.Errorf("上传进度状态非法：%v", s.UploadStatus.Status)
		}

	// ========================
	// Translate
	// ========================
	case WorkflowTranslateStart:
		switch s.TranslateStatus.Status {
		case WorkflowPending:
			s.TranslateStatus.Status = WorkflowOngoing
			s.TranslateStatus.StartedAt = &now
		case WorkflowOngoing:
			return fmt.Errorf("翻译进度已标记为进行中，请勿重复标记")
		case WorkflowCompleted:
			return fmt.Errorf("翻译进度已标记为已完成，无法更改")
		default:
			return fmt.Errorf("翻译进度状态非法：%v", s.TranslateStatus.Status)
		}

	case WorkflowTranslateComplete:
		switch s.TranslateStatus.Status {
		case WorkflowOngoing:
			s.TranslateStatus.Status = WorkflowCompleted
			s.TranslateStatus.CompletedAt = &now
		case WorkflowPending:
			return fmt.Errorf("翻译进度尚未标记为进行中，无法标记为完成")
		case WorkflowCompleted:
			return fmt.Errorf("翻译进度已标记为已完成，无法重复标记")
		default:
			return fmt.Errorf("翻译进度状态非法：%v", s.TranslateStatus.Status)
		}

	// ========================
	// Proofread
	// ========================
	case WorkflowProofreadStart:
		switch s.ProofreadStatus.Status {
		case WorkflowPending:
			s.ProofreadStatus.Status = WorkflowOngoing
			s.ProofreadStatus.StartedAt = &now
		case WorkflowOngoing:
			return fmt.Errorf("校对进度已标记为进行中，请勿重复标记")
		case WorkflowCompleted:
			return fmt.Errorf("校对进度已标记为已完成，无法更改")
		default:
			return fmt.Errorf("校对进度状态非法：%v", s.ProofreadStatus.Status)
		}

	case WorkflowProofreadComplete:
		switch s.ProofreadStatus.Status {
		case WorkflowOngoing:
			s.ProofreadStatus.Status = WorkflowCompleted
			s.ProofreadStatus.CompletedAt = &now
		case WorkflowPending:
			return fmt.Errorf("校对进度尚未标记为进行中，无法标记为完成")
		case WorkflowCompleted:
			return fmt.Errorf("校对进度已标记为已完成，无法重复标记")
		default:
			return fmt.Errorf("校对进度状态非法：%v", s.ProofreadStatus.Status)
		}

	// ========================
	// Typeset
	// ========================
	case WorkflowTypesetStart:
		switch s.TypesettStatus.Status {
		case WorkflowPending:
			s.TypesettStatus.Status = WorkflowOngoing
			s.TypesettStatus.StartedAt = &now
		case WorkflowOngoing:
			return fmt.Errorf("嵌字进度已标记为进行中，请勿重复标记")
		case WorkflowCompleted:
			return fmt.Errorf("嵌字进度已标记为已完成，无法更改")
		default:
			return fmt.Errorf("嵌字进度状态非法：%v", s.TypesettStatus.Status)
		}

	case WorkflowTypesetComplete:
		switch s.TypesettStatus.Status {
		case WorkflowOngoing:
			s.TypesettStatus.Status = WorkflowCompleted
			s.TypesettStatus.CompletedAt = &now
		case WorkflowPending:
			return fmt.Errorf("嵌字进度尚未标记为进行中，无法标记为完成")
		case WorkflowCompleted:
			return fmt.Errorf("嵌字进度已标记为已完成，无法重复标记")
		default:
			return fmt.Errorf("嵌字进度状态非法：%v", s.TypesettStatus.Status)
		}

	// ========================
	// Review
	// ========================
	case WorkflowReviewComplete:
		switch s.ReviewStatus.Status {
		case WorkflowOngoing:
			s.ReviewStatus.Status = WorkflowCompleted
			s.ReviewStatus.CompletedAt = &now
		case WorkflowPending:
			return fmt.Errorf("监修进度尚未标记为进行中，无法标记为完成")
		case WorkflowCompleted:
			return fmt.Errorf("监修进度已标记为已完成，无法重复标记")
		default:
			return fmt.Errorf("监修进度状态非法：%v", s.ReviewStatus.Status)
		}

	// ========================
	// Publish
	// ========================
	case WorkflowPublishComplete:
		switch s.PublishStatus.Status {
		case WorkflowOngoing:
			s.PublishStatus.Status = WorkflowCompleted
			s.PublishStatus.CompletedAt = &now
		case WorkflowPending:
			return fmt.Errorf("发布进度尚未标记为进行中，无法标记为完成")
		case WorkflowCompleted:
			return fmt.Errorf("发布进度已标记为已完成，无法重复标记")
		default:
			return fmt.Errorf("发布进度状态非法：%v", s.PublishStatus.Status)
		}

	default:
		return fmt.Errorf("未知的工作流转换：%s", t)
	}

	return nil
}

// 将各个进度的状态编码为一个 uint64，方便存储和查询
func (s *workflowState) Mask() uint64 {
	var mask uint64 = 0

	// 每个进度占用 2 位，状态值为 0、1、2 分别对应 00、01、10
	mask |= uint64(s.UploadStatus.Status) << uint64(WorkflowUpload)
	mask |= uint64(s.TranslateStatus.Status) << uint64(WorkflowTranslate)
	mask |= uint64(s.ProofreadStatus.Status) << uint64(WorkflowProofread)
	mask |= uint64(s.TypesettStatus.Status) << uint64(WorkflowTypesett)
	mask |= uint64(s.ReviewStatus.Status) << uint64(WorkflowReview)
	mask |= uint64(s.PublishStatus.Status) << uint64(WorkflowPublish)

	return mask
}

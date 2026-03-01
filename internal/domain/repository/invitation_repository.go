package repository

import "labelplus-next-web-be/internal/domain/model"

type InvitationRepository interface {
	List(executor Executor, options ...QueryOption) ([]model.InvitationInfo, error)
	GetInfoByInviteeQQ(executor Executor, inviteeQQ string) (*model.InvitationInfo, error)
	Save(executor Executor, invitation *model.InvitationCreation) (string, error) // 如果成功，返回主键 ID
	UpdateByID(executor Executor, invitationID string, update *model.InvitationPatch) error
	Delete(executor Executor, invitationID string) error
}

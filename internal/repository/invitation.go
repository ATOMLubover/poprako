package repository

import (
	"labelplus-next-web-be/internal/domain/model"
	intf "labelplus-next-web-be/internal/domain/repository"
	"labelplus-next-web-be/internal/repository/entity"
	"labelplus-next-web-be/internal/util"
)

type invitationRepository struct {
	executor intf.Executor
}

func NewInvitationRepository(executor intf.Executor) intf.InvitationRepository {
	return &invitationRepository{executor: executor}
}

func (r *invitationRepository) withTransaction(executor intf.Executor) intf.Executor {
	if executor != nil {
		return executor
	}

	return r.executor
}

func (r *invitationRepository) BeginTransaction() intf.Executor {
	return r.executor.Begin()
}

func (r *invitationRepository) List(executor intf.Executor, options ...intf.QueryOption) ([]model.InvitationInfo, error) {
	executor = r.withTransaction(executor)

	executor = executor.Table(entity.InvitationTable)
	for _, opt := range options {
		executor = opt(executor)
	}

	var rows []entity.InvitationRow

	if err := executor.Find(&rows).Error; err != nil {
		return nil, err
	}

	result := make([]model.InvitationInfo, len(rows))
	for i, row := range rows {
		result[i] = entity.ToInvitationInfo(row)
	}
	return result, nil
}

func (r *invitationRepository) GetByID(executor intf.Executor, invitationID string) (*model.InvitationInfo, error) {
	executor = r.withTransaction(executor)

	var row entity.InvitationRow
	if err := executor.
		Table(entity.InvitationTable).
		Where("id = ?", invitationID).
		First(&row).Error; err != nil {
		return nil, err
	}
	info := entity.ToInvitationInfo(row)
	return &info, nil
}

func (r *invitationRepository) GetByInviteeQQ(executor intf.Executor, inviteeQQ string) (*model.InvitationInfo, error) {
	executor = r.withTransaction(executor)

	var row entity.InvitationRow
	if err := executor.
		Table(entity.InvitationTable).
		Where("invitee_qq = ? AND pending = TRUE", inviteeQQ).
		First(&row).Error; err != nil {
		return nil, err
	}
	info := entity.ToInvitationInfo(row)
	return &info, nil
}

func (r *invitationRepository) Create(executor intf.Executor, creation *model.InvitationCreation) (string, error) {
	executor = r.withTransaction(executor)

	row := entity.InvitationRow{
		ID:              util.GenerateUUID(),
		InvitorID:       creation.InvitorID,
		TargetTeamID:    creation.TargetTeamID,
		InviteeQQ:       creation.InviteeQQ,
		InvitationCode:  creation.InvitationCode,
		ToBeRawProvider: creation.ToBeRawProvider,
		ToBeTranslator:  creation.ToBeTranslator,
		ToBeProofreader: creation.ToBeProofreader,
		ToBeTypesetter:  creation.ToBeTypesetter,
		ToBeReviewer:    creation.ToBeReviewer,
		ToBeUploader:    creation.ToBeUploader,
		ToBeAdmin:       creation.ToBeAdmin,
		Pending:         true,
	}
	if err := executor.Create(&row).Error; err != nil {
		return "", err
	}
	return row.ID, nil
}

func (r *invitationRepository) Update(executor intf.Executor, update *model.InvitationUpdate) error {
	executor = r.withTransaction(executor)

	return executor.
		Table(entity.InvitationTable).
		Where("id = ?", update.ID).
		Updates(map[string]any{
			"to_be_raw_provider": update.ToBeRawProvider,
			"to_be_translator":   update.ToBeTranslator,
			"to_be_proofreader":  update.ToBeProofreader,
			"to_be_typesetter":   update.ToBeTypesetter,
			"to_be_reviewer":     update.ToBeReviewer,
			"to_be_uploader":     update.ToBeUploader,
			"to_be_admin":        update.ToBeAdmin,
		}).Error
}

func (r *invitationRepository) Invalidate(executor intf.Executor, invitationID string) error {
	executor = r.withTransaction(executor)

	return executor.
		Table(entity.InvitationTable).
		Where("id = ?", invitationID).
		Update("pending", false).Error
}

func (r *invitationRepository) DeleteByID(executor intf.Executor, invitationID string) error {
	executor = r.withTransaction(executor)

	return executor.
		Where("id = ?", invitationID).
		Delete(&entity.InvitationRow{}).Error
}

package repository

import (
	"time"

	"labelplus-next-web-be/internal/domain/model"
	intf "labelplus-next-web-be/internal/domain/repository"
	"labelplus-next-web-be/internal/repository/entity"
	"labelplus-next-web-be/internal/util"
)

type memberRepository struct {
	executor intf.Executor
}

func NewMemberRepository(executor intf.Executor) intf.MemberRepository {
	return &memberRepository{executor: executor}
}

func (r *memberRepository) withTransaction(executor intf.Executor) intf.Executor {
	if executor != nil {
		return executor
	}
	return r.executor
}

func (r *memberRepository) BeginTransaction() intf.Executor {
	return r.executor.Begin()
}

func (r *memberRepository) List(executor intf.Executor, options ...intf.QueryOption) ([]model.MemberProfile, error) {
	executor = r.withTransaction(executor)

	db := executor.Model(&entity.MemberRow{}).Where("member_table.deleted_at IS NULL")
	for _, opt := range options {
		db = opt(db)
	}

	var rows []entity.MemberRow
	if err := db.Find(&rows).Error; err != nil {
		return nil, err
	}

	result := make([]model.MemberProfile, len(rows))
	for i, row := range rows {
		result[i] = entity.ToMemberProfile(row, nil)
	}
	return result, nil
}

func (r *memberRepository) ListWithUserInfo(executor intf.Executor, options ...intf.QueryOption) ([]model.MemberProfile, error) {
	executor = r.withTransaction(executor)

	db := executor.Table("member_table").
		Select(`member_table.*,
			user_table.name           AS user_name,
			user_table.qq             AS user_qq,
			user_table.avatar_url     AS user_avatar_url,
			user_table.is_super_admin AS user_is_super_admin,
			user_table.created_at     AS user_created_at,
			user_table.updated_at     AS user_updated_at`).
		Joins("LEFT JOIN user_table ON user_table.id = member_table.user_id AND user_table.deleted_at IS NULL").
		Where("member_table.deleted_at IS NULL")
	for _, opt := range options {
		db = opt(db)
	}

	var rows []entity.MemberWithUserRow

	if err := db.Find(&rows).Error; err != nil {
		return nil, err
	}

	result := make([]model.MemberProfile, len(rows))
	for i, row := range rows {
		userInfo := &model.UserInfo{
			ID:           row.UserID,
			Name:         row.UserName,
			QQ:           row.UserQQ,
			AvatarURL:    row.UserAvatarURL,
			IsSuperAdmin: row.UserIsSuperAdmin,
			CreatedAt:    row.UserCreatedAt,
			UpdatedAt:    row.UserUpdatedAt,
		}
		result[i] = entity.ToMemberProfile(row.MemberRow, userInfo)
	}
	return result, nil
}

func (r *memberRepository) Exist(executor intf.Executor, options ...intf.QueryOption) (bool, error) {
	executor = r.withTransaction(executor)

	db := executor.Model(&entity.MemberRow{}).Where("deleted_at IS NULL")
	for _, opt := range options {
		db = opt(db)
	}

	var count int64
	if err := db.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *memberRepository) GetByID(executor intf.Executor, memberID string) (*model.MemberProfile, error) {
	executor = r.withTransaction(executor)

	var row entity.MemberRow
	err := executor.
		Model(&entity.MemberRow{}).
		Where("id = ? AND deleted_at IS NULL", memberID).
		First(&row).Error
	if err != nil {
		return nil, err
	}
	profile := entity.ToMemberProfile(row, nil)
	return &profile, nil
}

func (r *memberRepository) Create(executor intf.Executor, creation *model.MemberCreation) (string, error) {
	executor = r.withTransaction(executor)

	now := time.Now()
	row := entity.MemberRow{
		ID:     util.GenerateUUID(),
		UserID: creation.UserID,
		TeamID: creation.TeamID,
	}
	if creation.ToBeRawProvider {
		row.AssignedRawProviderAt = &now
	}
	if creation.ToBeTranslator {
		row.AssignedTranslatorAt = &now
	}
	if creation.ToBeProofreader {
		row.AssignedProofreaderAt = &now
	}
	if creation.ToBeTypesetter {
		row.AssignedTypesetterAt = &now
	}
	if creation.ToBeReviewer {
		row.AssignedReviewerAt = &now
	}
	if creation.ToBeUploader {
		row.AssignedUploaderAt = &now
	}
	if creation.ToBeAdmin {
		row.AssignedAdminAt = &now
	}

	if err := executor.Create(&row).Error; err != nil {
		return "", err
	}
	return row.ID, nil
}

func (r *memberRepository) Update(executor intf.Executor, update *model.MemberUpdate) error {
	executor = r.withTransaction(executor)

	updates := map[string]any{}
	if update.AssignRawProvider.State() == util.OptionSome {
		v := update.AssignRawProvider.Unwrap()
		updates["assigned_raw_provider_at"] = &v
	}
	if update.AssignTranslator.State() == util.OptionSome {
		v := update.AssignTranslator.Unwrap()
		updates["assigned_translator_at"] = &v
	}
	if update.AssignProofreader.State() == util.OptionSome {
		v := update.AssignProofreader.Unwrap()
		updates["assigned_proofreader_at"] = &v
	}
	if update.AssignTypesetter.State() == util.OptionSome {
		v := update.AssignTypesetter.Unwrap()
		updates["assigned_typesetter_at"] = &v
	}
	if update.AssignReviewer.State() == util.OptionSome {
		v := update.AssignReviewer.Unwrap()
		updates["assigned_reviewer_at"] = &v
	}
	if update.AssignUploader.State() == util.OptionSome {
		v := update.AssignUploader.Unwrap()
		updates["assigned_uploader_at"] = &v
	}
	if update.AssignAdmin.State() == util.OptionSome {
		v := update.AssignAdmin.Unwrap()
		updates["assigned_admin_at"] = &v
	}

	if len(updates) == 0 {
		return nil
	}

	return executor.
		Model(&entity.MemberRow{}).
		Where("id = ?", update.ID).
		Updates(updates).Error
}

func (r *memberRepository) DeleteByID(executor intf.Executor, memberID string) error {
	executor = r.withTransaction(executor)

	return executor.
		Model(&entity.MemberRow{}).
		Where("id = ?", memberID).
		Update("deleted_at", time.Now()).Error
}

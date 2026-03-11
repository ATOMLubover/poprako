package repository

import (
	"time"

	"labelplus-next-web-be/internal/domain/model"
	intf "labelplus-next-web-be/internal/domain/repository"
	"labelplus-next-web-be/internal/infrastructure/repository/entity"
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

func (r *memberRepository) ListProfiles(executor intf.Executor, options ...intf.QueryOption) ([]model.MemberWithUserInfo, error) {
	executor = r.withTransaction(executor)

	executor = executor.Table("member_table").
		Select(`member_table.*,
			user_table.name           AS user_name,
			user_table.qq             AS user_qq,
			user_table.avatar_oss_key AS user_avatar_oss_key,
			user_table.is_avatar_uploaded AS user_is_avatar_uploaded,
			user_table.is_super_admin AS user_is_super_admin,
			user_table.created_at     AS user_created_at,
			user_table.updated_at     AS user_updated_at`).
		Joins("LEFT JOIN user_table ON user_table.id = member_table.user_id AND user_table.deleted_at IS NULL").
		Where("member_table.deleted_at IS NULL")
	for _, opt := range options {
		executor = opt(executor)
	}

	var rows []entity.MemberWithUserRow
	if err := executor.Find(&rows).Error; err != nil {
		return nil, err
	}

	result := make([]model.MemberWithUserInfo, len(rows))
	for i, row := range rows {
		userInfo := model.UserInfo{
			ID:               row.UserID,
			Name:             row.UserName,
			QQ:               row.UserQQ,
			AvatarOSSKey:     row.UserAvatarOSSKey,
			IsAvatarUploaded: row.UserIsAvatarUploaded,
			IsSuperAdmin:     row.UserIsSuperAdmin,
			CreatedAt:        row.UserCreatedAt,
			UpdatedAt:        row.UserUpdatedAt,
		}
		result[i] = entity.ToMemberProfile(row.MemberProfileRow, userInfo)
	}

	return result, nil
}

func (r *memberRepository) ListWithTeamInfo(executor intf.Executor, options ...intf.QueryOption) ([]model.MemberWithTeamInfo, error) {
	executor = r.withTransaction(executor)

	executor = executor.Table(entity.MemberTable).
		Select(`member_table.*,
			team_table.name              AS team_name,
			team_table.description       AS team_description,
			team_table.avatar_oss_key    AS team_avatar_oss_key,
			team_table.is_avatar_uploaded AS team_is_avatar_uploaded,
			team_table.created_at        AS team_created_at,
			team_table.updated_at        AS team_updated_at`).
		Joins("LEFT JOIN team_table ON team_table.id = member_table.team_id AND team_table.deleted_at IS NULL").
		Where("member_table.deleted_at IS NULL")

	for _, opt := range options {
		executor = opt(executor)
	}

	var rows []entity.MemberWithTeamRow
	if err := executor.Find(&rows).Error; err != nil {
		return nil, err
	}

	result := make([]model.MemberWithTeamInfo, len(rows))
	for i, row := range rows {
		result[i] = entity.ToMemberWithTeamInfo(row)
	}

	return result, nil
}

func (r *memberRepository) ListProfilesWithUserInfo(executor intf.Executor, options ...intf.QueryOption) ([]model.MemberWithUserInfo, error) {
	executor = r.withTransaction(executor)

	executor = executor.Table("member_table").
		Select(`member_table.*,
			user_table.name           AS user_name,
			user_table.qq             AS user_qq,
			user_table.avatar_oss_key AS user_avatar_oss_key,
			user_table.is_avatar_uploaded AS user_is_avatar_uploaded,
			user_table.is_super_admin AS user_is_super_admin,
			user_table.created_at     AS user_created_at,
			user_table.updated_at     AS user_updated_at`).
		Joins("LEFT JOIN user_table ON user_table.id = member_table.user_id AND user_table.deleted_at IS NULL").
		Where("member_table.deleted_at IS NULL")
	for _, opt := range options {
		executor = opt(executor)
	}

	var rows []entity.MemberWithUserRow

	if err := executor.Find(&rows).Error; err != nil {
		return nil, err
	}

	result := make([]model.MemberWithUserInfo, len(rows))
	for i, row := range rows {
		userInfo := model.UserInfo{
			ID:               row.UserID,
			Name:             row.UserName,
			QQ:               row.UserQQ,
			AvatarOSSKey:     row.UserAvatarOSSKey,
			IsAvatarUploaded: row.UserIsAvatarUploaded,
			IsSuperAdmin:     row.UserIsSuperAdmin,
			CreatedAt:        row.UserCreatedAt,
			UpdatedAt:        row.UserUpdatedAt,
		}
		result[i] = entity.ToMemberProfile(row.MemberProfileRow, userInfo)
	}
	return result, nil
}

func (r *memberRepository) Exist(executor intf.Executor, options ...intf.QueryOption) (bool, error) {
	executor = r.withTransaction(executor)

	executor = executor.Table(entity.MemberTable).
		Where("member_table.deleted_at IS NULL")
	for _, opt := range options {
		executor = opt(executor)
	}

	var count int64

	if err := executor.Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *memberRepository) Get(executor intf.Executor, options ...intf.QueryOption) (model.MemberInfo, error) {
	executor = r.withTransaction(executor)
	executor = executor.Table(entity.MemberTable).Where("member_table.deleted_at IS NULL")
	for _, opt := range options {
		executor = opt(executor)
	}

	var row entity.MemberProfileRow
	if err := executor.First(&row).Error; err != nil {
		return model.MemberInfo{}, err
	}

	return model.MemberInfo{
		ID:                row.ID,
		UserID:            row.UserID,
		AssignRawProvider: row.AssignedRawProviderAt,
		AssignTranslator:  row.AssignedTranslatorAt,
		AssignProofreader: row.AssignedProofreaderAt,
		AssignTypesetter:  row.AssignedTypesetterAt,
		AssignReviewer:    row.AssignedReviewerAt,
		AssignPublisher:   row.AssignedPublisherAt,
		AssignAdmin:       row.AssignedAdminAt,
	}, nil
}

func (r *memberRepository) GetProfile(executor intf.Executor, options ...intf.QueryOption) (model.MemberWithUserInfo, error) {
	executor = r.withTransaction(executor)
	executor = executor.Table(entity.MemberTable).Where("member_table.deleted_at IS NULL")
	for _, opt := range options {
		executor = opt(executor)
	}

	var row entity.MemberProfileRow
	if err := executor.First(&row).Error; err != nil {
		return model.MemberWithUserInfo{}, err
	}

	return entity.ToMemberProfile(row, model.UserInfo{}), nil
}

func (r *memberRepository) Create(executor intf.Executor, creation model.MemberCreation) (string, error) {
	executor = r.withTransaction(executor)

	now := time.Now()

	row := entity.MemberInsertRow{
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
	if creation.ToBePublisher {
		row.AssignedPublisherAt = &now
	}
	if creation.ToBeAdmin {
		row.AssignedAdminAt = &now
	}

	if err := executor.Create(&row).Error; err != nil {
		return "", err
	}

	return row.ID, nil
}

func (r *memberRepository) Update(executor intf.Executor, update model.MemberUpdate) error {
	executor = r.withTransaction(executor)

	updates := map[string]any{
		"assigned_raw_provider_at": update.AssignRawProvider,
		"assigned_translator_at":   update.AssignTranslator,
		"assigned_proofreader_at":  update.AssignProofreader,
		"assigned_typesetter_at":   update.AssignTypesetter,
		"assigned_reviewer_at":     update.AssignReviewer,
		"assigned_publisher_at":    update.AssignPublisher,
		"assigned_admin_at":        update.AssignAdmin,
	}

	return executor.
		Table(entity.MemberTable).
		Where("id = ?", update.ID).
		Updates(updates).Error
}

func (r *memberRepository) Delete(executor intf.Executor, memberID string) error {
	executor = r.withTransaction(executor)

	return executor.
		Table(entity.MemberTable).
		Where("id = ?", memberID).
		Update("deleted_at", time.Now()).Error
}

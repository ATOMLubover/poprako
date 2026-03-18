package repository

import (
	"time"

	"labelplus-next-web-be/internal/domain/model"
	intf "labelplus-next-web-be/internal/domain/repository"
	"labelplus-next-web-be/internal/infrastructure/repository/entity"
	"labelplus-next-web-be/internal/util"

	"gorm.io/gorm"
)

type userRepository struct {
	executor intf.Executor
}

func NewUserRepository(executor intf.Executor) intf.UserRepository {
	return &userRepository{executor: executor}
}

func (r *userRepository) withTransaction(executor intf.Executor) intf.Executor {
	if executor != nil {
		return executor
	}

	return r.executor
}

func (r *userRepository) BeginTransaction() intf.Executor {
	return r.executor.Begin()
}

func (r *userRepository) List(executor intf.Executor, options ...intf.QueryOption) ([]model.UserInfo, error) {
	executor = r.withTransaction(executor)

	executor = executor.Table(entity.UserTable).Where("deleted_at IS NULL")
	for _, opt := range options {
		executor = opt(executor)
	}

	var rows []entity.UserInfoRow

	if err := executor.Find(&rows).Error; err != nil {
		return nil, err
	}

	result := make([]model.UserInfo, len(rows))
	for i, row := range rows {
		result[i] = entity.ToUserInfo(row)
	}

	return result, nil
}

func (r *userRepository) Get(executor intf.Executor, options ...intf.QueryOption) (model.UserInfo, error) {
	executor = r.withTransaction(executor)

	executor = executor.Table(entity.UserTable).Where("deleted_at IS NULL")
	for _, opt := range options {
		executor = opt(executor)
	}

	var row entity.UserInfoRow

	if err := executor.First(&row).Error; err != nil {
		return model.UserInfo{}, err
	}

	info := entity.ToUserInfo(row)

	return info, nil
}

func (r *userRepository) GetStats(executor intf.Executor, options ...intf.QueryOption) (model.UserStats, error) {
	executor = r.withTransaction(executor)

	executor = executor.Table(entity.UserStatsTable)
	for _, option := range options {
		executor = option(executor)
	}

	var row entity.UserStatsRow

	if err := executor.First(&row).Error; err != nil {
		return model.UserStats{}, err
	}

	userStats := entity.ToUserStats(row)

	return userStats, nil
}

func (r *userRepository) CreateStats(executor intf.Executor, creation model.UserStatsCreation) error {
	executor = r.withTransaction(executor)

	row := entity.UserStatsRow{
		ID:                      util.GenerateUUID(),
		UserID:                  creation.UserID,
		TotalAssignmentCount:    creation.TotalAssignmentCount,
		ActiveAssignmentCount:   creation.ActiveAssignmentCount,
		FinishedAssignmentCount: creation.FinishedAssignmentCount,
	}

	return executor.Create(&row).Error
}

func (r *userRepository) IncrementStats(executor intf.Executor, delta model.UserStatsDelta) error {
	executor = r.withTransaction(executor)

	return executor.
		Table(entity.UserStatsTable).
		Where("user_id = ?", delta.UserID).
		Updates(map[string]any{
			"total_assignment_count":    gorm.Expr("total_assignment_count + ?", delta.TotalAssignmentCountDelta),
			"active_assignment_count":   gorm.Expr("active_assignment_count + ?", delta.ActiveAssignmentCountDelta),
			"finished_assignment_count": gorm.Expr("finished_assignment_count + ?", delta.FinishedAssignmentCountDelta),
		}).Error
}

func (r *userRepository) GetCredentials(executor intf.Executor, options ...intf.QueryOption) (model.UserCredentials, error) {
	executor = r.withTransaction(executor)

	executor = executor.Table(entity.UserTable).Where("deleted_at IS NULL")
	for _, opt := range options {
		executor = opt(executor)
	}

	var row entity.UserCredentialsRow

	if err := executor.First(&row).Error; err != nil {
		return model.UserCredentials{}, err
	}
	return model.NewUserCredentials(row.ID, row.PasswordHash), nil
}

func (r *userRepository) Create(executor intf.Executor, registration model.UserCreation) (string, error) {
	executor = r.withTransaction(executor)

	row := entity.UserInsertRow{
		ID:               util.GenerateUUID(),
		Name:             registration.Name,
		QQ:               registration.QQ,
		AvatarOSSKey:     "",
		IsAvatarUploaded: false,
		PasswordHash:     registration.PasswordHash,
		IsSuperAdmin:     false,
	}
	if err := executor.Create(&row).Error; err != nil {
		return "", err
	}

	return row.ID, nil
}

func (r *userRepository) ReserveAvatar(executor intf.Executor, id string, avatarOSSKey string) error {
	executor = r.withTransaction(executor)

	return executor.
		Table(entity.UserTable).
		Where("id = ? AND deleted_at IS NULL", id).
		Updates(map[string]any{
			"avatar_oss_key":     avatarOSSKey,
			"is_avatar_uploaded": false,
		}).Error
}

func (r *userRepository) ConfirmAvatarUploaded(executor intf.Executor, id string) error {
	executor = r.withTransaction(executor)

	return executor.
		Table(entity.UserTable).
		Where("id = ? AND deleted_at IS NULL AND avatar_oss_key <> ''", id).
		Update("is_avatar_uploaded", true).Error
}

func (r *userRepository) Update(executor intf.Executor, update model.UserUpdate) error {
	executor = r.withTransaction(executor)

	return executor.
		Table(entity.UserTable).
		Where("id = ? AND deleted_at IS NULL", update.ID).
		Updates(map[string]any{
			"name":          update.Name,
			"qq":            update.QQ,
			"password_hash": update.PasswordHash,
		}).Error
}

func (r *userRepository) Delete(executor intf.Executor, id string) error {
	executor = r.withTransaction(executor)

	return executor.
		Model(&entity.UserInfoRow{}).
		Where("id = ?", id).
		Update("deleted_at", time.Now()).Error
}

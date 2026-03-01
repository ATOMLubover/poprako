package application

import (
	"errors"
	"time"

	"labelplus-next-web-be/internal/domain/model"
	"labelplus-next-web-be/internal/domain/repository"
	"labelplus-next-web-be/internal/domain/service"
	"labelplus-next-web-be/internal/repository/query_option"
	"labelplus-next-web-be/internal/util"
	"labelplus-next-web-be/internal/value"

	"go.uber.org/zap"
)

type InvitationApplication interface {
	ListInvitations(
		scope util.TraceScope,
		currentUserID string,
	) ([]*value.InvitationInfo, error)
	CreateInvitation(
		scope util.TraceScope,
		currentUserID string,
		args *value.CreateInvitationArgs,
	) (*value.InvitationInfo, error)
	PatchInvitation(
		scope util.TraceScope,
		currentUserID string,
		args *value.PatchInvitationArgs,
	) error
	DeleteInvitation(
		scope util.TraceScope,
		currentUserID string,
		invitationID string,
	) error
}

type invitationApplication struct {
	userRepository       repository.UserRepository
	invitationRepository repository.InvitationRepository

	permissionService service.PermissionService
	invitationService service.InvitationService
}

func NewInvitationApplication(
	userRepository repository.UserRepository,
	invitationRepository repository.InvitationRepository,
	permissionService service.PermissionService,
	invitationService service.InvitationService,
) InvitationApplication {
	if userRepository == nil ||
		invitationRepository == nil ||
		permissionService == nil ||
		invitationService == nil {
		zap.L().Panic(
			"NewInvitationApplication: 依赖项不能为空",
			zap.Bool("userRepository_nil", userRepository == nil),
			zap.Bool("invitationRepository_nil", invitationRepository == nil),
			zap.Bool("permissionService_nil", permissionService == nil),
			zap.Bool("invitationService_nil", invitationService == nil),
		)
	}

	return &invitationApplication{
		userRepository:       userRepository,
		invitationRepository: invitationRepository,
		permissionService:    permissionService,
		invitationService:    invitationService,
	}
}

func (ia *invitationApplication) ListInvitations(
	scope util.TraceScope,
	currentUserID string,
) ([]*value.InvitationInfo, error) {
	const fn = "InvitationApplication.ListInvitations"

	if currentUserID == "" {
		scope.Logger().Warn(fn + ": currentUserID 为空")
		return nil, errors.New(ErrInternalError)
	}

	scope.WithField(
		zap.String("current_user_id", currentUserID),
	)

	scope.Logger().Debug(fn + ": 被调用")

	// 鉴权：检查当前用户是否有权限查看邀请信息
	userInfo, err := ia.userRepository.GetInfoByID(nil, currentUserID)
	if err != nil {
		scope.Logger().Error(fn+": 获取用户信息失败", zap.Error(err))
		return nil, errors.New("无法获取用户信息")
	}

	if !ia.permissionService.CheckPermission(userInfo, model.PermissionInvitationsList) {
		return nil, errors.New("没有权限查看邀请信息")
	}

	// 获取邀请信息列表
	invitationList, err := ia.invitationRepository.List(nil, query_option.CreatedAtDesc()...)
	if err != nil {
		scope.Logger().Error(fn+": 获取邀请信息列表失败", zap.Error(err))
		return nil, errors.New("无法获取邀请信息列表")
	}

	// 转换为应用层的值对象
	result := make([]*value.InvitationInfo, len(invitationList))
	for i, info := range invitationList {
		result[i] = value.NewInvitationInfoFromModel(&info)
	}

	return result, nil
}

func (ia *invitationApplication) CreateInvitation(
	scope util.TraceScope,
	currentUserID string,
	args *value.CreateInvitationArgs,
) (*value.InvitationInfo, error) {
	const fn = "InvitationApplication.CreateInvitation"

	if currentUserID == "" {
		scope.Logger().Warn(fn + ": currentUserID 为空")
		return nil, errors.New(ErrInternalError)
	}

	if args == nil {
		scope.Logger().Warn(fn + ": args 为空")
		return nil, errors.New(ErrInternalError)
	}

	if err := args.Validate(); err != nil {
		scope.Logger().Warn(fn+": 参数验证失败", zap.Error(err))
		return nil, errors.New("参数错误: " + err.Error())
	}

	scope.WithField(
		zap.String("current_user_id", currentUserID),
		zap.Any("args", args),
	)

	scope.Logger().Debug(fn + ": 被调用")

	// 鉴权：检查当前用户是否有权限创建邀请
	userInfo, err := ia.userRepository.GetInfoByID(nil, currentUserID)
	if err != nil {
		scope.Logger().Error(fn+": 获取用户信息失败", zap.Error(err))
		return nil, errors.New("无法获取用户信息")
	}

	if !ia.permissionService.CheckPermission(userInfo, model.PermissionInvitationsCreate) {
		return nil, errors.New("没有权限创建邀请")
	}

	// 检查是否已经有相关 QQ 的成员存在
	isUserExisting, err := ia.userRepository.ExistsByQQ(nil, args.InviteeQQ)
	if err != nil {
		scope.Logger().Error(fn+": 检查用户是否存在失败", zap.Error(err))
		return nil, errors.New("无法检查用户信息")
	}

	if isUserExisting {
		return nil, errors.New("已经有对应 QQ 的成员存在")
	}

	// 创建邀请信息
	invitationCode, err := ia.invitationService.GenerateInvitationCode()
	if err != nil {
		scope.Logger().Error(fn+": 生成邀请码失败", zap.Error(err))
		return nil, errors.New("生成邀请码失败")
	}

	invitationCreation := model.NewInvitationCreation(
		currentUserID,
		args.InviteeQQ,
		invitationCode,
		model.UnmaskRoles(args.Roles)...,
	)

	// 此处应当使用乐观锁进行更新，不采用事务
	// 如果有 pending 的同 QQ 邀请存在，则会直接覆盖，不会创建新的邀请记录
	invitationID, err := ia.invitationRepository.Save(nil, invitationCreation)
	if err != nil {
		scope.Logger().Error(fn+": 创建邀请信息失败", zap.Error(err))
		return nil, errors.New("创建邀请失败")
	}

	// 构造创建成功的邀请信息值对象
	invitationInfo := value.NewInvitationInfo(
		invitationID,
		currentUserID,
		args.InviteeQQ,
		true,
		args.Roles,
		time.Now().UnixMilli(),
	)

	return invitationInfo, nil
}

func (ia *invitationApplication) PatchInvitation(
	scope util.TraceScope,
	currentUserID string,
	args *value.PatchInvitationArgs,
) error {
	const fn = "InvitationApplication.PatchInvitation"

	if currentUserID == "" {
		scope.Logger().Warn(fn + ": currentUserID 为空")
		return errors.New(ErrInternalError)
	}

	if args == nil {
		scope.Logger().Warn(fn + ": args 为空")
		return errors.New(ErrInternalError)
	}

	if err := args.Validate(); err != nil {
		scope.Logger().Warn(fn+": 参数验证失败", zap.Error(err))
		return errors.New("参数错误: " + err.Error())
	}

	scope.WithField(
		zap.String("current_user_id", currentUserID),
		zap.Any("args", args),
	)

	scope.Logger().Debug(fn + ": 被调用")

	userInfo, err := ia.userRepository.GetInfoByID(nil, currentUserID)
	if err != nil {
		scope.Logger().Error(fn+": 获取用户信息失败", zap.Error(err))
		return errors.New("无法获取用户信息")
	}

	if !ia.permissionService.CheckPermission(userInfo, model.PermissionInvitationsPatch) {
		return errors.New("没有权限修改邀请信息")
	}

	invitationPatch := model.NewInvitationPatch(
		args.ID,
		model.UnmaskRoles(args.Roles)...,
	)

	err = ia.invitationRepository.UpdateByID(nil, args.ID, invitationPatch)
	if err != nil {
		scope.Logger().Error(fn+": 更新邀请信息失败", zap.Error(err))
		return errors.New("更新邀请失败")
	}

	return nil
}

func (ia *invitationApplication) DeleteInvitation(
	scope util.TraceScope,
	currentUserID string,
	invitationID string,
) error {
	const fn = "InvitationApplication.DeleteInvitation"

	scope.WithField(
		zap.String("current_user_id", currentUserID),
		zap.String("invitation_id", invitationID),
	)

	scope.Logger().Debug(fn + ": 被调用")

	// 鉴权：检查当前用户是否有权限删除邀请
	userInfo, err := ia.userRepository.GetInfoByID(nil, currentUserID)
	if err != nil {
		scope.Logger().Error(fn+": 获取用户信息失败", zap.Error(err))
		return errors.New("无法获取用户信息")
	}

	if !ia.permissionService.CheckPermission(userInfo, model.PermissionInvitationsDelete) {
		return errors.New("没有权限删除邀请")
	}

	// 删除邀请信息
	err = ia.invitationRepository.Delete(nil, invitationID)
	if err != nil {
		scope.Logger().Error(fn+": 删除邀请信息失败", zap.Error(err))
		return errors.New("删除邀请失败")
	}

	return nil
}

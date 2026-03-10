package model

import "go.uber.org/zap"

type (
	OnLoadUserInfo       func(userID string) (UserInfo, error)
	OnLoadMemberInfo     func(teamID string, userID string) (MemberInfo, error)
	OnLoadAssignmentInfo func(comicID string, userID string) (AssignmentInfo, error)
	OnLoadComicInfo      func(comicID string) (ComicInfo, error)
	OnLoadChapterInfo    func(chapterID string) (ChapterInfo, error)
	OnLoadPageInfo       func(pageID string) (PageInfo, error)
)

func isSuperAdmin(
	checkName string,
	userID string,
	onLoadUserInfo OnLoadUserInfo,
) bool {
	userInfo, ok := loadUserInfoForCheck(checkName, userID, onLoadUserInfo)
	if !ok {
		return false
	}

	return userInfo.IsSuperAdmin
}

func loadUserInfoForCheck(
	checkName string,
	userID string,
	onLoadUserInfo OnLoadUserInfo,
) (UserInfo, bool) {
	userInfo, err := onLoadUserInfo(userID)
	if err != nil {
		zap.L().Error(
			checkName+": 获取用户信息失败",
			zap.String("userID", userID),
			zap.Error(err),
		)

		return UserInfo{}, false
	}

	return userInfo, true
}

func loadMemberInfoForCheck(
	checkName string,
	userID string,
	teamID string,
	onLoadMemberInfo OnLoadMemberInfo,
) (MemberInfo, bool) {
	memberInfo, err := onLoadMemberInfo(teamID, userID)
	if err != nil {
		zap.L().Error(
			checkName+": 获取成员信息失败",
			zap.String("teamID", teamID),
			zap.String("userID", userID),
			zap.Error(err),
		)

		return MemberInfo{}, false
	}

	return memberInfo, true
}

func isTeamAdmin(
	checkName string,
	userID string,
	teamID string,
	onLoadMemberInfo OnLoadMemberInfo,
) bool {
	memberInfo, ok := loadMemberInfoForCheck(checkName, userID, teamID, onLoadMemberInfo)
	if !ok {
		return false
	}

	return memberInfo.HasAnyRole(RoleAdmin)
}

func loadComicInfoForCheck(
	checkName string,
	comicID string,
	onLoadComicInfo OnLoadComicInfo,
) (ComicInfo, bool) {
	comicInfo, err := onLoadComicInfo(comicID)
	if err != nil {
		zap.L().Error(
			checkName+": 获取漫画信息失败",
			zap.String("comicID", comicID),
			zap.Error(err),
		)

		return ComicInfo{}, false
	}

	return comicInfo, true
}

func loadAssignmentInfoForCheck(
	checkName string,
	chapterID string,
	userID string,
	onLoadAssignmentInfo OnLoadAssignmentInfo,
) (AssignmentInfo, bool) {
	assignmentInfo, err := onLoadAssignmentInfo(chapterID, userID)
	if err != nil {
		zap.L().Error(
			checkName+": 获取分配信息失败",
			zap.String("comicID", chapterID),
			zap.String("userID", userID),
			zap.Error(err),
		)

		return AssignmentInfo{}, false
	}

	return assignmentInfo, true
}

func loadChapterInfoForCheck(
	checkName string,
	chapterID string,
	onLoadChapterInfo OnLoadChapterInfo,
) (ChapterInfo, bool) {
	chapterInfo, err := onLoadChapterInfo(chapterID)
	if err != nil {
		zap.L().Error(
			checkName+": 获取章节信息失败",
			zap.String("chapterID", chapterID),
			zap.Error(err),
		)

		return ChapterInfo{}, false
	}

	return chapterInfo, true
}

func loadPageInfoForCheck(
	checkName string,
	pageID string,
	onLoadPageInfo OnLoadPageInfo,
) (PageInfo, bool) {
	pageInfo, err := onLoadPageInfo(pageID)
	if err != nil {
		zap.L().Error(
			checkName+": 获取页面信息失败",
			zap.String("pageID", pageID),
			zap.Error(err),
		)

		return PageInfo{}, false
	}

	return pageInfo, true
}

// 类型化所有的权限，实现更加模块化和可定制化的权限判断

type (
	permInvitationList   struct{}
	permInvitationCreate struct{}
	permInvitationDelete struct{}
	permInvitationUpdate struct{}
)

func PermInvitationList() permInvitationList     { return permInvitationList{} }
func PermInvitationCreate() permInvitationCreate { return permInvitationCreate{} }
func PermInvitationDelete() permInvitationDelete { return permInvitationDelete{} }
func PermInvitationUpdate() permInvitationUpdate { return permInvitationUpdate{} }

func (permInvitationList) Check(
	userID string,
	teamID string,
	onLoadMemberInfo OnLoadMemberInfo,
) bool {
	// 只有汉化组管理员可以查看邀请列表
	return isTeamAdmin("PermInvitationList.Check", userID, teamID, onLoadMemberInfo)
}

func (permInvitationCreate) Check(
	userID string,
	teamID string,
	onLoadMemberInfo OnLoadMemberInfo,
) bool {
	// 只有汉化组管理员可以创建邀请
	return isTeamAdmin("PermInvitationCreate.Check", userID, teamID, onLoadMemberInfo)
}

func (permInvitationDelete) Check(
	userID string,
	teamID string,
	onLoadMemberInfo OnLoadMemberInfo,
) bool {
	// 只有汉化组管理员可以删除邀请
	return isTeamAdmin("PermInvitationDelete.Check", userID, teamID, onLoadMemberInfo)
}

func (permInvitationUpdate) Check(
	userID string,
	teamID string,
	onLoadMemberInfo OnLoadMemberInfo,
) bool {
	// 只有汉化组管理员可以更新邀请
	return isTeamAdmin("PermInvitationUpdate.Check", userID, teamID, onLoadMemberInfo)
}

type (
	permUserList   struct{}
	permUserView   struct{}
	permUserUpdate struct{}
	permUserRemove struct{} // remove 是硬删除，而 delete 是软删除
)

func PermUserList() permUserList     { return permUserList{} }
func PermUserView() permUserView     { return permUserView{} }
func PermUserUpdate() permUserUpdate { return permUserUpdate{} }
func PermUserRemove() permUserRemove { return permUserRemove{} }

func (permUserList) Check(
	userID string,
	onLoadUserInfo OnLoadUserInfo,
) bool {
	return isSuperAdmin("PermUserList.Check", userID, onLoadUserInfo)
}

func (permUserView) Check(
	userID string,
	targetUserID string,
) bool {
	// FIXME: 应该只有当 user 与 targetUser 在至少同一个汉化组中时
	// 才能查看 targetUser 的信息，否则只能查看自己的信息
	// 当前暂时不设置权限管理
	return true
}

func (permUserUpdate) Check(
	userID string,
	targetUserID string,
	onLoadUserInfo OnLoadUserInfo,
) bool {
	if userID == targetUserID {
		return true
	}

	return isSuperAdmin("PermUserUpdate.Check", userID, onLoadUserInfo)
}

func (permUserRemove) Check(
	userID string,
	targetUserID string,
	onLoadUserInfo OnLoadUserInfo,
) bool {
	if !isSuperAdmin("PermUserRemove.Check", userID, onLoadUserInfo) {
		return false
	}

	// 仅超级管理员有权限删除用户，且不能删除自己
	return userID != targetUserID
}

type (
	permTeamListAll struct{}
	permTeamCreate  struct{}
	permTeamUpdate  struct{}
	permTeamRemove  struct{}
)

func PermTeamListAll() permTeamListAll { return permTeamListAll{} }
func PermTeamCreate() permTeamCreate   { return permTeamCreate{} }
func PermTeamUpdate() permTeamUpdate   { return permTeamUpdate{} }
func PermTeamRemove() permTeamRemove   { return permTeamRemove{} }

func (permTeamListAll) Check(
	userID string,
	onLoadUserInfo OnLoadUserInfo,
) bool {
	return isSuperAdmin("PermTeamListAll.Check", userID, onLoadUserInfo)
}

func (permTeamCreate) Check(
	userID string,
	onLoadUserInfo OnLoadUserInfo,
) bool {
	return isSuperAdmin("PermTeamCreate.Check", userID, onLoadUserInfo)
}

func (permTeamUpdate) Check(
	userID string,
	teamID string,
	onLoadUserInfo OnLoadUserInfo,
	onLoadMemberInfo OnLoadMemberInfo,
) bool {
	if isSuperAdmin("PermTeamUpdate.Check", userID, onLoadUserInfo) {
		return true
	}

	// 汉化组管理员也可以更新所属汉化组
	return isTeamAdmin("PermTeamUpdate.Check", userID, teamID, onLoadMemberInfo)
}

func (permTeamRemove) Check(
	userID string,
	teamID string,
	onLoadUserInfo OnLoadUserInfo,
	onLoadMemberInfo OnLoadMemberInfo,
) bool {
	return isSuperAdmin("PermTeamRemove.Check", userID, onLoadUserInfo)
}

type (
	permMemberList   struct{}
	permMemberCreate struct{}
	permMemberUpdate struct{}
	permMemberRemove struct{}
)

func PermMemberList() permMemberList     { return permMemberList{} }
func PermMemberCreate() permMemberCreate { return permMemberCreate{} }
func PermMemberUpdate() permMemberUpdate { return permMemberUpdate{} }
func PermMemberRemove() permMemberRemove { return permMemberRemove{} }

func (permMemberCreate) Check(
	userID string,
	onLoadUserInfo OnLoadUserInfo,
) bool {
	return isSuperAdmin("PermMemberCreate.Check", userID, onLoadUserInfo)
}

func (permMemberList) Check(
	userID string,
	teamID string,
	onLoadMemberInfo OnLoadMemberInfo,
) bool {
	// 仅管理员可以查看成员列表
	return isTeamAdmin("PermMemberList.Check", userID, teamID, onLoadMemberInfo)
}

func (permMemberUpdate) Check(
	userID string,
	teamID string,
	onLoadMemberInfo OnLoadMemberInfo,
) bool {
	// 仅管理员可以更新成员
	return isTeamAdmin("PermMemberUpdate.Check", userID, teamID, onLoadMemberInfo)
}

func (permMemberRemove) Check(
	userID string,
	teamID string,
	onLoadMemberInfo OnLoadMemberInfo,
) bool {
	// 仅管理员可以删除成员
	return isTeamAdmin("PermMemberRemove.Check", userID, teamID, onLoadMemberInfo)
}

type (
	permComicList   struct{}
	permComicCreate struct{}
	permComicUpdate struct{}
	permComicDelete struct{}
)

func PermComicList() permComicList     { return permComicList{} }
func PermComicCreate() permComicCreate { return permComicCreate{} }
func PermComicUpdate() permComicUpdate { return permComicUpdate{} }
func PermComicDelete() permComicDelete { return permComicDelete{} }

func (permComicList) Check(
	userID string,
	teamID string,
	onLoadMemberInfo OnLoadMemberInfo,
) bool {
	// 只要是汉化组成员，就可以查看所属汉化组的漫画列表
	_, ok := loadMemberInfoForCheck("PermComicList.Check", userID, teamID, onLoadMemberInfo)
	if !ok {
		return false
	}

	// HACK: 理论上只要返回了成员信息，就说明是成员了
	return true
}

func (permComicCreate) Check(
	userID string,
	teamID string,
	onLoadMemberInfo OnLoadMemberInfo,
) bool {
	// 只有汉化组管理员可以创建漫画
	return isTeamAdmin("PermComicCreate.Check", userID, teamID, onLoadMemberInfo)
}

func (permComicUpdate) Check(
	userID string,
	teamID string,
	onLoadMemberInfo OnLoadMemberInfo,
) bool {
	// 只有汉化组管理员可以更新漫画
	return isTeamAdmin("PermComicUpdate.Check", userID, teamID, onLoadMemberInfo)
}

func (permComicDelete) Check(
	userID string,
	teamID string,
	onLoadMemberInfo OnLoadMemberInfo,
) bool {
	// 只有汉化组管理员可以删除漫画
	return isTeamAdmin("PermComicDelete.Check", userID, teamID, onLoadMemberInfo)
}

type (
	permChapterList   struct{}
	permChapterCreate struct{}
	permChapterUpdate struct{}
	permChapterDelete struct{}
)

func PermChapterList() permChapterList     { return permChapterList{} }
func PermChapterCreate() permChapterCreate { return permChapterCreate{} }
func PermChapterUpdate() permChapterUpdate { return permChapterUpdate{} }
func PermChapterDelete() permChapterDelete { return permChapterDelete{} }

func (permChapterList) Check(
	userID string,
	comicID string,
	onLoadMemberInfo OnLoadMemberInfo,
	onLoadComicInfo OnLoadComicInfo,
) bool {
	// 先查找漫画信息，然后根据其所属汉化组来判断用户是否有权限查看章节列表
	comicInfo, ok := loadComicInfoForCheck("PermChapterList.Check", comicID, onLoadComicInfo)
	if !ok {
		return false
	}

	_, ok = loadMemberInfoForCheck("PermChapterList.Check", userID, comicInfo.TeamID, onLoadMemberInfo)
	if !ok {
		return false
	}

	// HACK: 只要是所属汉化组的成员，就可以查看章节列表
	return true
}

func (permChapterCreate) Check(
	userID string,
	comicID string,
	onLoadMemberInfo OnLoadMemberInfo,
	onLoadComicInfo OnLoadComicInfo,
) bool {
	// 先查找漫画信息，然后根据其所属汉化组来判断用户是否有权限创建章节
	comicInfo, ok := loadComicInfoForCheck("PermChapterCreate.Check", comicID, onLoadComicInfo)
	if !ok {
		return false
	}

	return isTeamAdmin("PermChapterCreate.Check", userID, comicInfo.TeamID, onLoadMemberInfo)
}

func (permChapterUpdate) Check(
	userID string,
	comicID string,
	onLoadMemberInfo OnLoadMemberInfo,
	onLoadComicInfo OnLoadComicInfo,
) bool {
	// 先查找漫画信息，然后根据其所属汉化组来判断用户是否有权限更新章节
	comicInfo, ok := loadComicInfoForCheck("PermChapterUpdate.Check", comicID, onLoadComicInfo)
	if !ok {
		return false
	}

	return isTeamAdmin("PermChapterUpdate.Check", userID, comicInfo.TeamID, onLoadMemberInfo)
}

func (permChapterDelete) Check(
	userID string,
	comicID string,
	onLoadMemberInfo OnLoadMemberInfo,
	onLoadComicInfo OnLoadComicInfo,
) bool {
	// 先查找漫画信息，然后根据其所属汉化组来判断用户是否有权限删除章节
	comicInfo, ok := loadComicInfoForCheck("PermChapterDelete.Check", comicID, onLoadComicInfo)
	if !ok {
		return false
	}

	return isTeamAdmin("PermChapterDelete.Check", userID, comicInfo.TeamID, onLoadMemberInfo)
}

type (
	permAssignmentList   struct{}
	permAssignmentCreate struct{}
	permAssignmentUpdate struct{}
	permAssignmentDelete struct{}
)

func PermAssignmentList() permAssignmentList     { return permAssignmentList{} }
func PermAssignmentCreate() permAssignmentCreate { return permAssignmentCreate{} }
func PermAssignmentUpdate() permAssignmentUpdate { return permAssignmentUpdate{} }
func PermAssignmentDelete() permAssignmentDelete { return permAssignmentDelete{} }

func (permAssignmentList) Check(
	userID string,
	chapterID string,
	onLoadChapterInfo OnLoadChapterInfo,
	onLoadComicInfo OnLoadComicInfo,
	onLoadMemberInfo OnLoadMemberInfo,
) bool {
	// 先查找章节信息，然后根据其所属漫画的所属汉化组来判断用户是否有权限查看分配列表
	chapterInfo, ok := loadChapterInfoForCheck("PermAssignmentList.Check", chapterID, onLoadChapterInfo)
	if !ok {
		return false
	}

	comicInfo, ok := loadComicInfoForCheck("PermAssignmentList.Check", chapterInfo.ComicID, onLoadComicInfo)
	if !ok {
		return false
	}

	_, ok = loadMemberInfoForCheck("PermAssignmentList.Check", userID, comicInfo.TeamID, onLoadMemberInfo)
	if !ok {
		return false
	}

	// HACK: 只要是所属汉化组的成员，就可以查看分配列表
	return true
}

func (permAssignmentCreate) Check(
	userID string,
	chapterID string,
	onLoadAssignmentInfo OnLoadAssignmentInfo,
) bool {
	assignmentInfo, ok := loadAssignmentInfoForCheck("PermAssignmentCreate.Check", chapterID, userID, onLoadAssignmentInfo)
	if !ok {
		return false
	}

	// 只有当用户是 reviewer 时，才有权限创建分配
	return assignmentInfo.HasAnyRole(RoleReviewer)
}

func (permAssignmentUpdate) Check(
	userID string,
	chapterID string,
	onLoadAssignmentInfo OnLoadAssignmentInfo,
) bool {
	assignmentInfo, ok := loadAssignmentInfoForCheck("PermAssignmentUpdate.Check", chapterID, userID, onLoadAssignmentInfo)
	if !ok {
		return false
	}

	// 只有当用户是 reviewer 时，才有权限更新分配
	return assignmentInfo.HasAnyRole(RoleReviewer)
}

func (permAssignmentDelete) Check(
	userID string,
	chapterID string,
	onLoadAssignmentInfo OnLoadAssignmentInfo,
) bool {
	assignmentInfo, ok := loadAssignmentInfoForCheck("PermAssignmentDelete.Check", chapterID, userID, onLoadAssignmentInfo)
	if !ok {
		return false
	}

	// 只有当用户是 reviewer 时，才有权限删除分配
	return assignmentInfo.HasAnyRole(RoleReviewer)
}

type (
	permPageList   struct{}
	permPageCreate struct{}
	permPageUpdate struct{}
	permPageDelete struct{}
)

func PermPageList() permPageList     { return permPageList{} }
func PermPageCreate() permPageCreate { return permPageCreate{} }
func PermPageUpdate() permPageUpdate { return permPageUpdate{} }
func PermPageDelete() permPageDelete { return permPageDelete{} }

func (permPageList) Check(
	userID string,
	chapterID string,
	onLoadChapterInfo OnLoadChapterInfo,
	onLoadComicInfo OnLoadComicInfo,
	onLoadMemberInfo OnLoadMemberInfo,
) bool {
	// 先查找章节信息，然后根据其所属漫画的所属汉化组来判断用户是否有权限查看页面列表
	chapterInfo, ok := loadChapterInfoForCheck("PermPageList.Check", chapterID, onLoadChapterInfo)
	if !ok {
		return false
	}

	comicInfo, ok := loadComicInfoForCheck("PermPageList.Check", chapterInfo.ComicID, onLoadComicInfo)
	if !ok {
		return false
	}

	_, ok = loadMemberInfoForCheck("PermPageList.Check", userID, comicInfo.TeamID, onLoadMemberInfo)
	if !ok {
		return false
	}

	// HACK: 只要是所属汉化组的成员，就可以查看页面列表
	return true
}

func (permPageCreate) Check(
	userID string,
	chapterID string,
	onLoadAssignmentInfo OnLoadAssignmentInfo,
) bool {
	assignmentInfo, ok := loadAssignmentInfoForCheck("PermPageCreate.Check", chapterID, userID, onLoadAssignmentInfo)
	if !ok {
		return false
	}

	// 只有当用户有 reviewer 或 raw_provider 分工时，才有权限创建页面
	return assignmentInfo.HasAnyRole(RoleReviewer, RoleRawProvider)
}

func (permPageUpdate) Check(
	userID string,
	pageID string,
	onLoadPageInfo OnLoadPageInfo,
	onLoadChapterInfo OnLoadChapterInfo,
	onLoadAssignmentInfo OnLoadAssignmentInfo,
) bool {
	pageInfo, ok := loadPageInfoForCheck("PermPageUpdate.Check", pageID, onLoadPageInfo)
	if !ok {
		return false
	}

	chapterInfo, ok := loadChapterInfoForCheck("PermPageUpdate.Check", pageInfo.ChapterID, onLoadChapterInfo)
	if !ok {
		return false
	}

	assignmentInfo, ok := loadAssignmentInfoForCheck("PermPageUpdate.Check", chapterInfo.ComicID, userID, onLoadAssignmentInfo)
	if !ok {
		return false
	}

	// 只有当用户有 reviewer 或 raw_provider 分工时，才有权限更新页面
	return assignmentInfo.HasAnyRole(RoleReviewer, RoleRawProvider)
}

func (permPageDelete) Check(
	userID string,
	pageID string,
	onLoadPageInfo OnLoadPageInfo,
	onLoadChapterInfo OnLoadChapterInfo,
	onLoadAssignmentInfo OnLoadAssignmentInfo,
) bool {
	pageInfo, ok := loadPageInfoForCheck("PermPageDelete.Check", pageID, onLoadPageInfo)
	if !ok {
		return false
	}

	chapterInfo, ok := loadChapterInfoForCheck("PermPageDelete.Check", pageInfo.ChapterID, onLoadChapterInfo)
	if !ok {
		return false
	}

	assignmentInfo, ok := loadAssignmentInfoForCheck("PermPageDelete.Check", chapterInfo.ComicID, userID, onLoadAssignmentInfo)
	if !ok {
		return false
	}

	// 只有当用户有 reviewer 或 raw_provider 分工时，才有权限删除页面
	return assignmentInfo.HasAnyRole(RoleReviewer, RoleRawProvider)
}

type permUnitSave struct {
	role RoleFlag
}

func PermUnitSave(role RoleFlag) permUnitSave {
	return permUnitSave{role: role}
}

func (p permUnitSave) Check(
	userID string,
	pageID string,
	onLoadPageInfo OnLoadPageInfo,
	onLoadChapterInfo OnLoadChapterInfo,
	onLoadAssignmentInfo OnLoadAssignmentInfo,
) bool {
	// 先查找页面信息，然后根据其所属章节的分工信息来判断用户是否有权限保存单元
	pageInfo, ok := loadPageInfoForCheck("PermUnitSave.Check", pageID, onLoadPageInfo)
	if !ok {
		return false
	}

	chapterInfo, ok := loadChapterInfoForCheck("PermUnitSave.Check", pageInfo.ChapterID, onLoadChapterInfo)
	if !ok {
		return false
	}

	assignmentInfo, ok := loadAssignmentInfoForCheck("PermUnitSave.Check", chapterInfo.ComicID, userID, onLoadAssignmentInfo)
	if !ok {
		return false
	}

	// 只有当用户有对应分工时，才有权限保存单元
	return assignmentInfo.HasAnyRole(p.role)
}

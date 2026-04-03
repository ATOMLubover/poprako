package app

import (
	"context"
	"testing"

	"poprako-s/internal/app/val"
)

type assignmentAppStub struct {
	listByChapterCalled bool
	listMyCalled        bool
	createCalled        bool
	updateCalled        bool
	removeCalled        bool
}

func (s *assignmentAppStub) ListByChapter(context.Context, string, *val.ListChapterAssignmentArgs) ([]*val.AssignmentInfo, error) {
	s.listByChapterCalled = true
	return []*val.AssignmentInfo{{ID: "assignment-1"}}, nil
}

func (s *assignmentAppStub) ListMy(context.Context, string, *val.ListMyAssignmentArgs) ([]*val.AssignmentInfo, error) {
	s.listMyCalled = true
	return []*val.AssignmentInfo{{ID: "assignment-1"}}, nil
}

func (s *assignmentAppStub) Create(context.Context, string, *val.CreateAssignmentArgs) (*val.CreateAssignmentRes, error) {
	s.createCalled = true
	return &val.CreateAssignmentRes{ID: "assignment-1"}, nil
}

func (s *assignmentAppStub) Update(context.Context, string, *val.UpdateAssignmentArgs) error {
	s.updateCalled = true
	return nil
}

func (s *assignmentAppStub) Remove(context.Context, string, string) error {
	s.removeCalled = true
	return nil
}

type chapterAppStub struct {
	listCalled   bool
	createCalled bool
	updateCalled bool
	removeCalled bool
}

func (s *chapterAppStub) List(context.Context, string, *val.ListChapterArgs) ([]*val.ChapterInfo, error) {
	s.listCalled = true
	return []*val.ChapterInfo{{ID: "chapter-1"}}, nil
}

func (s *chapterAppStub) Create(context.Context, string, *val.CreateChapterArgs) (*val.CreateChapterRes, error) {
	s.createCalled = true
	return &val.CreateChapterRes{ID: "chapter-1"}, nil
}

func (s *chapterAppStub) Update(context.Context, string, *val.UpdateChapterArgs) error {
	s.updateCalled = true
	return nil
}

func (s *chapterAppStub) Remove(context.Context, string, string) error {
	s.removeCalled = true
	return nil
}

type comicAppStub struct {
	listCalled   bool
	createCalled bool
	updateCalled bool
	removeCalled bool
}

func (s *comicAppStub) List(context.Context, string, *val.ListComicArgs) ([]*val.ComicInfo, error) {
	s.listCalled = true
	return []*val.ComicInfo{{ID: "comic-1"}}, nil
}

func (s *comicAppStub) Create(context.Context, string, *val.CreateComicArgs) (*val.CreateComicRes, error) {
	s.createCalled = true
	return &val.CreateComicRes{ID: "comic-1"}, nil
}

func (s *comicAppStub) Update(context.Context, string, *val.UpdateComicArgs) error {
	s.updateCalled = true
	return nil
}

func (s *comicAppStub) Remove(context.Context, string, string) error {
	s.removeCalled = true
	return nil
}

type invitationAppStub struct {
	listCalled   bool
	createCalled bool
	updateCalled bool
	removeCalled bool
}

func (s *invitationAppStub) List(context.Context, string, *val.ListTeamInvitationArgs) ([]*val.InvitationInfo, error) {
	s.listCalled = true
	return []*val.InvitationInfo{{ID: "inv-1"}}, nil
}

func (s *invitationAppStub) Create(context.Context, string, *val.CreateInvitationArgs) (*val.InvitationInfo, error) {
	s.createCalled = true
	return &val.InvitationInfo{ID: "inv-1"}, nil
}

func (s *invitationAppStub) Update(context.Context, string, *val.UpdateInvitationArgs) error {
	s.updateCalled = true
	return nil
}

func (s *invitationAppStub) Remove(context.Context, string, string) error {
	s.removeCalled = true
	return nil
}

type memberAppStub struct {
	createCalled     bool
	listByTeamCalled bool
	listMyCalled     bool
	updateCalled     bool
	removeCalled     bool
	joinTeamCalled   bool
}

func (s *memberAppStub) Create(context.Context, string, *val.CreateMemberArgs) (*val.CreateMemberRes, error) {
	s.createCalled = true
	return &val.CreateMemberRes{ID: "member-1"}, nil
}

func (s *memberAppStub) ListByTeam(context.Context, string, *val.ListTeamMemberArgs) ([]*val.MemberInfo, error) {
	s.listByTeamCalled = true
	return []*val.MemberInfo{{ID: "member-1"}}, nil
}

func (s *memberAppStub) ListMy(context.Context, string, *val.ListMyMemberArgs) ([]*val.MemberInfo, error) {
	s.listMyCalled = true
	return []*val.MemberInfo{{ID: "member-1"}}, nil
}

func (s *memberAppStub) UpdateRole(context.Context, string, *val.UpdateMemberRoleArgs) error {
	s.updateCalled = true
	return nil
}

func (s *memberAppStub) Remove(context.Context, string, string) error {
	s.removeCalled = true
	return nil
}

func (s *memberAppStub) JoinTeam(context.Context, string, *val.JoinTeamArgs) error {
	s.joinTeamCalled = true
	return nil
}

type pageAppStub struct {
	reserveCalled bool
	listCalled    bool
	updateCalled  bool
	removeCalled  bool
}

func (s *pageAppStub) Reserve(context.Context, string, *val.ReserveChapterPagesArgs) (*val.ReserveChapterPagesRes, error) {
	s.reserveCalled = true
	return &val.ReserveChapterPagesRes{}, nil
}

func (s *pageAppStub) List(context.Context, string, *val.ListChapterPageArgs) ([]*val.PageInfo, error) {
	s.listCalled = true
	return []*val.PageInfo{{ID: "page-1"}}, nil
}

func (s *pageAppStub) Update(context.Context, string, *val.UpdatePageArgs) error {
	s.updateCalled = true
	return nil
}

func (s *pageAppStub) Remove(context.Context, string, string) error {
	s.removeCalled = true
	return nil
}

type teamAppStub struct {
	createCalled  bool
	listCalled    bool
	listMyCalled  bool
	updateCalled  bool
	removeCalled  bool
	reserveCalled bool
	confirmCalled bool
}

func (s *teamAppStub) Create(context.Context, string, *val.CreateTeamArgs) (*val.CreateTeamRes, error) {
	s.createCalled = true
	return &val.CreateTeamRes{ID: "team-1"}, nil
}

func (s *teamAppStub) List(context.Context, string, *val.ListTeamArgs) ([]*val.TeamInfo, error) {
	s.listCalled = true
	return []*val.TeamInfo{{ID: "team-1"}}, nil
}

func (s *teamAppStub) ListMy(context.Context, string, *val.ListMyTeamArgs) ([]*val.TeamInfo, error) {
	s.listMyCalled = true
	return []*val.TeamInfo{{ID: "team-1"}}, nil
}

func (s *teamAppStub) Update(context.Context, string, *val.UpdateTeamArgs) error {
	s.updateCalled = true
	return nil
}

func (s *teamAppStub) Remove(context.Context, string, string) error {
	s.removeCalled = true
	return nil
}

func (s *teamAppStub) ReserveAvatar(context.Context, string, string) (*val.ReserveTeamAvatarRes, error) {
	s.reserveCalled = true
	return &val.ReserveTeamAvatarRes{PutURL: "put-url"}, nil
}

func (s *teamAppStub) ConfirmAvatarUploaded(context.Context, string, string) error {
	s.confirmCalled = true
	return nil
}

type unitAppStub struct {
	listCalled bool
	saveCalled bool
}

func (s *unitAppStub) List(context.Context, string, string) ([]*val.UnitInfo, error) {
	s.listCalled = true
	return []*val.UnitInfo{{ID: "unit-1"}}, nil
}

func (s *unitAppStub) Save(context.Context, string, *val.SavePageUnitArgs) error {
	s.saveCalled = true
	return nil
}

type userAppStub struct {
	parseTokenCalled bool
	loginCalled      bool
	regCalled        bool
	getCalled        bool
	getMyCalled      bool
	updateCalled     bool
	statsCalled      bool
	reserveCalled    bool
	confirmCalled    bool
	removeCalled     bool
}

func (s *userAppStub) ParseToken(context.Context, string) (string, error) {
	s.parseTokenCalled = true
	return "user-1", nil
}

func (s *userAppStub) Login(context.Context, *val.LoginUserArgs) (*val.LoginUserRes, error) {
	s.loginCalled = true
	return &val.LoginUserRes{UserID: "user-1", AccessToken: "token"}, nil
}

func (s *userAppStub) Reg(context.Context, *val.RegUserArgs) (*val.RegUserRes, error) {
	s.regCalled = true
	return &val.RegUserRes{UserID: "user-1", AccessToken: "token"}, nil
}

func (s *userAppStub) GetInfo(context.Context, string) (*val.UserInfo, error) {
	s.getCalled = true
	return &val.UserInfo{ID: "user-1"}, nil
}

func (s *userAppStub) GetMyInfo(context.Context, string) (*val.UserInfo, error) {
	s.getMyCalled = true
	return &val.UserInfo{ID: "user-1"}, nil
}

func (s *userAppStub) UpdateMyInfo(context.Context, string, *val.UpdateUserArgs) error {
	s.updateCalled = true
	return nil
}

func (s *userAppStub) GetMyStats(context.Context, string) (*val.UserStatsInfo, error) {
	s.statsCalled = true
	return &val.UserStatsInfo{UserID: "user-1"}, nil
}

func (s *userAppStub) ReserveMyAvatar(context.Context, string) (*val.ReserveUserAvatarRes, error) {
	s.reserveCalled = true
	return &val.ReserveUserAvatarRes{PutURL: "put-url"}, nil
}

func (s *userAppStub) ConfirmMyAvatarUploaded(context.Context, string) error {
	s.confirmCalled = true
	return nil
}

func (s *userAppStub) Remove(context.Context, string, string) error {
	s.removeCalled = true
	return nil
}

type worksetAppStub struct {
	listCalled   bool
	createCalled bool
	updateCalled bool
	removeCalled bool
}

func (s *worksetAppStub) List(context.Context, string, *val.ListWorksetArgs) ([]*val.WorksetInfo, error) {
	s.listCalled = true
	return []*val.WorksetInfo{{ID: "workset-1"}}, nil
}

func (s *worksetAppStub) Create(context.Context, string, *val.CreateWorksetArgs) (*val.CreateWorksetRes, error) {
	s.createCalled = true
	return &val.CreateWorksetRes{ID: "workset-1"}, nil
}

func (s *worksetAppStub) Update(context.Context, string, *val.UpdateWorksetArgs) error {
	s.updateCalled = true
	return nil
}

func (s *worksetAppStub) Remove(context.Context, string, string) error {
	s.removeCalled = true
	return nil
}

func assertPanics(t *testing.T, fn func()) {
	t.Helper()

	defer func() {

		if recover() == nil {
			t.Fatal("expected panic")
		}
	}()
	fn()
}

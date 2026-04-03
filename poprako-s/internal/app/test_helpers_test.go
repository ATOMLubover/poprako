package app

import (
	"context"
	"errors"
	"testing"
	"time"

	"poprako-s/internal/domain/model"
	repoiface "poprako-s/internal/domain/repo"
	mock_event "poprako-s/internal/infra/event/mock"
	mock_oss "poprako-s/internal/infra/oss/mock"
	mock_repo "poprako-s/internal/infra/repo/mock"
)

func ptr[T any](v T) *T {
	return &v
}

func ptrptr[T any](v T) **T {
	p := &v
	return &p
}

func adminMember() *model.MemberInfo {
	now := time.Now()
	return &model.MemberInfo{ID: "member-1", UserID: "user-1", TeamID: "team-1", AssignedAdminAt: &now}
}

func reviewerMember() *model.MemberInfo {
	now := time.Now()
	return &model.MemberInfo{ID: "member-1", UserID: "user-1", TeamID: "team-1", AssignedReviewerAt: &now}
}

func reviewerAssignment() *model.AssignmentInfo {
	now := time.Now()
	return &model.AssignmentInfo{ID: "assignment-1", ChapterID: "chapter-1", UserID: "user-1", AssignedReviewerAt: &now}
}

func rawProviderAssignment() *model.AssignmentInfo {
	now := time.Now()
	return &model.AssignmentInfo{ID: "assignment-1", ChapterID: "chapter-1", UserID: "user-1", AssignedRawProviderAt: &now}
}

func translatorAssignment() *model.AssignmentInfo {
	now := time.Now()
	return &model.AssignmentInfo{ID: "assignment-1", ChapterID: "chapter-1", UserID: "user-1", AssignedTranslatorAt: &now}
}

func proofreaderAssignment() *model.AssignmentInfo {
	now := time.Now()
	return &model.AssignmentInfo{ID: "assignment-1", ChapterID: "chapter-1", UserID: "user-1", AssignedProofreaderAt: &now}
}

func publisherAssignment() *model.AssignmentInfo {
	now := time.Now()
	return &model.AssignmentInfo{ID: "assignment-1", ChapterID: "chapter-1", UserID: "user-1", AssignedPublisherAt: &now}
}

func superAdminUser() *model.UserInfo {
	now := time.Now()
	return &model.UserInfo{ID: "user-1", QQ: "10001", Name: "super", IsSuperAdmin: true, LastLoginAt: now, CreatedAt: now, UpdatedAt: now}
}

func normalUser() *model.UserInfo {
	now := time.Now()
	return &model.UserInfo{ID: "user-1", QQ: "10001", Name: "normal", LastLoginAt: now, CreatedAt: now, UpdatedAt: now}
}

func background() context.Context {
	return context.Background()
}

func newMockOSSClient() *mock_oss.Client {
	client, ok := mock_oss.NewMockOSSClient().(*mock_oss.Client)
	if !ok {
		panic("mock_oss.NewMockOSSClient did not return *mock_oss.Client")
	}
	return client
}

func newMockEventBus() *mock_event.EventBus {
	bus, ok := mock_event.NewMockEventBus().(*mock_event.EventBus)
	if !ok {
		panic("mock_event.NewMockEventBus did not return *mock_event.EventBus")
	}
	return bus
}

type mockTxnRepos struct {
	assignment *mock_repo.AssignmentRepo
	chapter    *mock_repo.ChapterRepo
	comic      *mock_repo.ComicRepo
	invitation *mock_repo.InvitationRepo
	member     *mock_repo.MemberRepo
	page       *mock_repo.PageRepo
	team       *mock_repo.TeamRepo
	unit       *mock_repo.UnitRepo
	user       *mock_repo.UserRepo
	workset    *mock_repo.WorksetRepo
}

func newMockTxnContext(repos mockTxnRepos) context.Context {
	cx := context.Background()
	if repos.assignment != nil {
		cx = mock_repo.WithMockAssignmentRepo(cx, repos.assignment)
	}
	if repos.chapter != nil {
		cx = mock_repo.WithMockChapterRepo(cx, repos.chapter)
	}
	if repos.comic != nil {
		cx = mock_repo.WithMockComicRepo(cx, repos.comic)
	}
	if repos.invitation != nil {
		cx = mock_repo.WithMockInvitationRepo(cx, repos.invitation)
	}
	if repos.member != nil {
		cx = mock_repo.WithMockMemberRepo(cx, repos.member)
	}
	if repos.page != nil {
		cx = mock_repo.WithMockPageRepo(cx, repos.page)
	}
	if repos.team != nil {
		cx = mock_repo.WithMockTeamRepo(cx, repos.team)
	}
	if repos.unit != nil {
		cx = mock_repo.WithMockUnitRepo(cx, repos.unit)
	}
	if repos.user != nil {
		cx = mock_repo.WithMockUserRepo(cx, repos.user)
	}
	if repos.workset != nil {
		cx = mock_repo.WithMockWorksetRepo(cx, repos.workset)
	}
	return cx
}

func requireNoErr(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

type erroringAssignmentRepo struct {
	repoiface.AssignmentRepo
	getByIDErr error
	getErr     error
	listErr    error
	updateErr  error
	deleteErr  error
	fromTxnErr error
}

func (r *erroringAssignmentRepo) GetByID(id string) (*model.AssignmentInfo, error) {
	if r.getByIDErr != nil {
		return nil, r.getByIDErr
	}
	return r.AssignmentRepo.GetByID(id)
}

func (r *erroringAssignmentRepo) Get(opt model.AssignmentQueryOpt) (*model.AssignmentInfo, error) {
	if r.getErr != nil {
		return nil, r.getErr
	}
	return r.AssignmentRepo.Get(opt)
}

func (r *erroringAssignmentRepo) List(opt model.AssignmentQueryOpt) ([]model.AssignmentInfo, error) {
	if r.listErr != nil {
		return nil, r.listErr
	}
	return r.AssignmentRepo.List(opt)
}

func (r *erroringAssignmentRepo) Update(u *model.AssignmentUpdate) error {
	if r.updateErr != nil {
		return r.updateErr
	}
	return r.AssignmentRepo.Update(u)
}

func (r *erroringAssignmentRepo) Delete(id string) error {
	if r.deleteErr != nil {
		return r.deleteErr
	}
	return r.AssignmentRepo.Delete(id)
}

func (r *erroringAssignmentRepo) FromTxnCx(cx context.Context) (repoiface.AssignmentRepo, error) {
	if r.fromTxnErr != nil {
		return nil, r.fromTxnErr
	}
	return r.AssignmentRepo.FromTxnCx(cx)
}

type erroringUserRepo struct {
	repoiface.UserRepo
	getCredsErr error
	getByIDErr  error
	getByQQErr  error
	updateErr   error
	statsErr    error
	preFillErr  error
	confirmErr  error
	removeErr   error
	fromTxnErr  error
}

func (r *erroringUserRepo) GetCredsByQQ(qq string) (*model.UserCreds, error) {
	if r.getCredsErr != nil {
		return nil, r.getCredsErr
	}
	return r.UserRepo.GetCredsByQQ(qq)
}

func (r *erroringUserRepo) GetByID(id string) (*model.UserInfo, error) {
	if r.getByIDErr != nil {
		return nil, r.getByIDErr
	}
	return r.UserRepo.GetByID(id)
}

func (r *erroringUserRepo) GetByQQ(qq string) (*model.UserInfo, error) {
	if r.getByQQErr != nil {
		return nil, r.getByQQErr
	}
	return r.UserRepo.GetByQQ(qq)
}

func (r *erroringUserRepo) Update(u *model.UserUpdate) error {
	if r.updateErr != nil {
		return r.updateErr
	}
	return r.UserRepo.Update(u)
}

func (r *erroringUserRepo) GetOrCreateStats(userID string) (*model.UserStats, error) {
	if r.statsErr != nil {
		return nil, r.statsErr
	}
	return r.UserRepo.GetOrCreateStats(userID)
}

func (r *erroringUserRepo) PreFillAvatarOSSKey(id string, avatarOSSKey string) error {
	if r.preFillErr != nil {
		return r.preFillErr
	}
	return r.UserRepo.PreFillAvatarOSSKey(id, avatarOSSKey)
}

func (r *erroringUserRepo) ConfirmAvatarUploaded(id string) error {
	if r.confirmErr != nil {
		return r.confirmErr
	}
	return r.UserRepo.ConfirmAvatarUploaded(id)
}

func (r *erroringUserRepo) Remove(id string) error {
	if r.removeErr != nil {
		return r.removeErr
	}
	return r.UserRepo.Remove(id)
}

func (r *erroringUserRepo) FromTxnCx(cx context.Context) (repoiface.UserRepo, error) {
	if r.fromTxnErr != nil {
		return nil, r.fromTxnErr
	}
	return r.UserRepo.FromTxnCx(cx)
}

type erroringMemberRepo struct {
	repoiface.MemberRepo
	getByIDErr error
	getErr     error
	listErr    error
	existErr   error
	updateErr  error
	deleteErr  error
	fromTxnErr error
}

func (r *erroringMemberRepo) GetByID(id string) (*model.MemberInfo, error) {
	if r.getByIDErr != nil {
		return nil, r.getByIDErr
	}
	return r.MemberRepo.GetByID(id)
}

func (r *erroringMemberRepo) Get(opt model.MemberQueryOpt) (*model.MemberInfo, error) {
	if r.getErr != nil {
		return nil, r.getErr
	}
	return r.MemberRepo.Get(opt)
}

func (r *erroringMemberRepo) List(opt model.MemberQueryOpt) ([]model.MemberInfo, error) {
	if r.listErr != nil {
		return nil, r.listErr
	}
	return r.MemberRepo.List(opt)
}

func (r *erroringMemberRepo) Exist(opt model.MemberQueryOpt) (bool, error) {
	if r.existErr != nil {
		return false, r.existErr
	}
	return r.MemberRepo.Exist(opt)
}

func (r *erroringMemberRepo) Update(u *model.MemberUpdate) error {
	if r.updateErr != nil {
		return r.updateErr
	}
	return r.MemberRepo.Update(u)
}

func (r *erroringMemberRepo) Delete(id string) error {
	if r.deleteErr != nil {
		return r.deleteErr
	}
	return r.MemberRepo.Delete(id)
}

func (r *erroringMemberRepo) FromTxnCx(cx context.Context) (repoiface.MemberRepo, error) {
	if r.fromTxnErr != nil {
		return nil, r.fromTxnErr
	}
	return r.MemberRepo.FromTxnCx(cx)
}

type erroringPageRepo struct {
	repoiface.PageRepo
	getByIDErr    error
	listErr       error
	createBatchErr error
	updateErr     error
	deleteErr     error
	fromTxnErr    error
}

func (r *erroringPageRepo) GetByID(id string) (*model.PageInfo, error) {
	if r.getByIDErr != nil {
		return nil, r.getByIDErr
	}
	return r.PageRepo.GetByID(id)
}

func (r *erroringPageRepo) List(opt model.PageQueryOpt) ([]model.PageInfo, error) {
	if r.listErr != nil {
		return nil, r.listErr
	}
	return r.PageRepo.List(opt)
}

func (r *erroringPageRepo) CreateBatch(pages []*model.PageCreation) error {
	if r.createBatchErr != nil {
		return r.createBatchErr
	}
	return r.PageRepo.CreateBatch(pages)
}

func (r *erroringPageRepo) Update(u *model.PageUpdate) error {
	if r.updateErr != nil {
		return r.updateErr
	}
	return r.PageRepo.Update(u)
}

func (r *erroringPageRepo) Delete(id string) error {
	if r.deleteErr != nil {
		return r.deleteErr
	}
	return r.PageRepo.Delete(id)
}

func (r *erroringPageRepo) FromTxnCx(cx context.Context) (repoiface.PageRepo, error) {
	if r.fromTxnErr != nil {
		return nil, r.fromTxnErr
	}
	return r.PageRepo.FromTxnCx(cx)
}

type erroringChapterRepo struct {
	repoiface.ChapterRepo
	getByIDErr error
	updateErr  error
	removeErr  error
	fromTxnErr error
}

func (r *erroringChapterRepo) GetByID(id string) (*model.ChapterInfo, error) {
	if r.getByIDErr != nil {
		return nil, r.getByIDErr
	}
	return r.ChapterRepo.GetByID(id)
}

func (r *erroringChapterRepo) Update(u *model.ChapterUpdate) error {
	if r.updateErr != nil {
		return r.updateErr
	}
	return r.ChapterRepo.Update(u)
}

func (r *erroringChapterRepo) Remove(id string) error {
	if r.removeErr != nil {
		return r.removeErr
	}
	return r.ChapterRepo.Remove(id)
}

func (r *erroringChapterRepo) FromTxnCx(cx context.Context) (repoiface.ChapterRepo, error) {
	if r.fromTxnErr != nil {
		return nil, r.fromTxnErr
	}
	return r.ChapterRepo.FromTxnCx(cx)
}

type erroringComicRepo struct {
	repoiface.ComicRepo
	getByIDErr error
	updateErr  error
	deleteErr  error
	fromTxnErr error
}

func (r *erroringComicRepo) GetByID(id string) (*model.ComicInfo, error) {
	if r.getByIDErr != nil {
		return nil, r.getByIDErr
	}
	return r.ComicRepo.GetByID(id)
}

func (r *erroringComicRepo) Update(u *model.ComicUpdate) error {
	if r.updateErr != nil {
		return r.updateErr
	}
	return r.ComicRepo.Update(u)
}

func (r *erroringComicRepo) Delete(id string) error {
	if r.deleteErr != nil {
		return r.deleteErr
	}
	return r.ComicRepo.Delete(id)
}

func (r *erroringComicRepo) FromTxnCx(cx context.Context) (repoiface.ComicRepo, error) {
	if r.fromTxnErr != nil {
		return nil, r.fromTxnErr
	}
	return r.ComicRepo.FromTxnCx(cx)
}

type erroringWorksetRepo struct {
	repoiface.WorksetRepo
	getByIDErr error
	updateErr  error
	deleteErr  error
	fromTxnErr error
}

func (r *erroringWorksetRepo) GetByID(id string) (*model.WorksetInfo, error) {
	if r.getByIDErr != nil {
		return nil, r.getByIDErr
	}
	return r.WorksetRepo.GetByID(id)
}

func (r *erroringWorksetRepo) Update(u *model.WorksetUpdate) error {
	if r.updateErr != nil {
		return r.updateErr
	}
	return r.WorksetRepo.Update(u)
}

func (r *erroringWorksetRepo) Delete(id string) error {
	if r.deleteErr != nil {
		return r.deleteErr
	}
	return r.WorksetRepo.Delete(id)
}

func (r *erroringWorksetRepo) FromTxnCx(cx context.Context) (repoiface.WorksetRepo, error) {
	if r.fromTxnErr != nil {
		return nil, r.fromTxnErr
	}
	return r.WorksetRepo.FromTxnCx(cx)
}

type erroringInvitationRepo struct {
	repoiface.InvitationRepo
	getByInviteeQQErr error
	listErr           error
	updateErr         error
	invalidateErr     error
	deleteErr         error
	fromTxnErr        error
}

func (r *erroringInvitationRepo) GetByInviteeQQ(qq string) (*model.InvitationInfo, error) {
	if r.getByInviteeQQErr != nil {
		return nil, r.getByInviteeQQErr
	}
	return r.InvitationRepo.GetByInviteeQQ(qq)
}

func (r *erroringInvitationRepo) List(opt model.InvitationQueryOpt) ([]model.InvitationInfo, error) {
	if r.listErr != nil {
		return nil, r.listErr
	}
	return r.InvitationRepo.List(opt)
}

func (r *erroringInvitationRepo) Update(u *model.InvitationUpdate) error {
	if r.updateErr != nil {
		return r.updateErr
	}
	return r.InvitationRepo.Update(u)
}

func (r *erroringInvitationRepo) Invalidate(id string) error {
	if r.invalidateErr != nil {
		return r.invalidateErr
	}
	return r.InvitationRepo.Invalidate(id)
}

func (r *erroringInvitationRepo) Delete(id string) error {
	if r.deleteErr != nil {
		return r.deleteErr
	}
	return r.InvitationRepo.Delete(id)
}

func (r *erroringInvitationRepo) FromTxnCx(cx context.Context) (repoiface.InvitationRepo, error) {
	if r.fromTxnErr != nil {
		return nil, r.fromTxnErr
	}
	return r.InvitationRepo.FromTxnCx(cx)
}

type erroringTeamRepo struct {
	repoiface.TeamRepo
	listErr    error
	updateErr  error
	deleteErr  error
	preFillErr error
	confirmErr error
}

func (r *erroringTeamRepo) List(opt model.TeamQueryOpt) ([]model.TeamInfo, error) {
	if r.listErr != nil {
		return nil, r.listErr
	}
	return r.TeamRepo.List(opt)
}

func (r *erroringTeamRepo) Update(u *model.TeamUpdate) error {
	if r.updateErr != nil {
		return r.updateErr
	}
	return r.TeamRepo.Update(u)
}

func (r *erroringTeamRepo) Delete(id string) error {
	if r.deleteErr != nil {
		return r.deleteErr
	}
	return r.TeamRepo.Delete(id)
}

func (r *erroringTeamRepo) PreFillAvatarOSSKey(id string, avatarOSSKey string) error {
	if r.preFillErr != nil {
		return r.preFillErr
	}
	return r.TeamRepo.PreFillAvatarOSSKey(id, avatarOSSKey)
}

func (r *erroringTeamRepo) ConfirmAvatarUploaded(id string) error {
	if r.confirmErr != nil {
		return r.confirmErr
	}
	return r.TeamRepo.ConfirmAvatarUploaded(id)
}

var errBoom = errors.New("boom")

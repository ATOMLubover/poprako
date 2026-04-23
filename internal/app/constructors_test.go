package app

import (
	"testing"

	"poprako-s/internal/cfg"
	"poprako-s/internal/domain/service"
	mock_repo "poprako-s/internal/infra/repo/mock"
)

func TestConstructorsRejectNilDependencies(t *testing.T) {
	assertPanics(t, func() {
		NewAssignmentApp(nil, mock_repo.NewMockAssignmentRepo(), mock_repo.NewMockChapterInvitationRepo(), mock_repo.NewMockChapterRepo(), mock_repo.NewMockComicRepo(), mock_repo.NewMockMemberRepo(), mock_repo.NewMockWorksetRepo(), mock_repo.NewMockUserRepo(), mock_repo.NewMockTxnMgr(nil), newMockEventBus(), newMockOSSClient())
	})
	assertPanics(t, func() {
		NewChapterApp(nil, service.NewChapterInvitationService(), mock_repo.NewMockMemberRepo(), mock_repo.NewMockWorksetRepo(), mock_repo.NewMockComicRepo(), mock_repo.NewMockChapterRepo(), mock_repo.NewMockAssignmentRepo(), mock_repo.NewMockUserRepo(), mock_repo.NewMockPageRepo(), mock_repo.NewMockChapterInvitationRepo(), mock_repo.NewMockTxnMgr(nil), mock_repo.NewMockOSSMessageRepo(), newMockEventBus(), newMockOSSClient())
	})
	assertPanics(t, func() {
		NewComicApp(nil, mock_repo.NewMockMemberRepo(), mock_repo.NewMockWorksetRepo(), mock_repo.NewMockComicRepo(), mock_repo.NewMockTxnMgr(nil), mock_repo.NewMockOSSMessageRepo(), newMockEventBus(), newMockOSSClient())
	})
	assertPanics(t, func() {
		NewInvitationApp(nil, mock_repo.NewMockMemberRepo(), mock_repo.NewMockInvitationRepo())
	})
	assertPanics(t, func() {
		NewMemberApp(nil, mock_repo.NewMockUserRepo(), mock_repo.NewMockMemberRepo(), mock_repo.NewMockTeamRepo(), mock_repo.NewMockInvitationRepo(), mock_repo.NewMockTxnMgr(nil), newMockOSSClient())
	})
	assertPanics(t, func() {
		NewPageApp(nil, mock_repo.NewMockAssignmentRepo(), mock_repo.NewMockChapterRepo(), mock_repo.NewMockPageRepo(), mock_repo.NewMockTxnMgr(nil), mock_repo.NewMockOSSMessageRepo(), newMockOSSClient())
	})
	assertPanics(t, func() {
		NewTeamApp(nil, service.NewMemberService(), mock_repo.NewMockUserRepo(), mock_repo.NewMockTeamRepo(), mock_repo.NewMockMemberRepo(), mock_repo.NewMockTxnMgr(nil), mock_repo.NewMockOSSMessageRepo(), newMockOSSClient())
	})
	assertPanics(t, func() {
		NewUnitApp(nil, newMockEventBus(), mock_repo.NewMockUserRepo(), mock_repo.NewMockPageRepo(), mock_repo.NewMockChapterRepo(), mock_repo.NewMockAssignmentRepo(), mock_repo.NewMockUnitRepo(), mock_repo.NewMockTxnMgr(nil))
	})
	assertPanics(t, func() {
		NewUserApp(nil, service.NewMemberService(), mock_repo.NewMockUserRepo(), mock_repo.NewMockInvitationRepo(), mock_repo.NewMockMemberRepo(), mock_repo.NewMockTxnMgr(nil), mock_repo.NewMockOSSMessageRepo(), newMockEventBus(), newMockOSSClient(), &cfg.AuthCfg{SecretKey: "secret", ExpHrs: 1})
	})
	assertPanics(t, func() {
		NewWorksetApp(nil, mock_repo.NewMockMemberRepo(), mock_repo.NewMockWorksetRepo(), mock_repo.NewMockTxnMgr(nil))
	})
}

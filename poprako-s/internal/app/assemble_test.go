package app

import (
	"errors"
	"testing"
	"time"

	"poprako-s/internal/domain/model"
)

func TestAssembleHelpers(t *testing.T) {
	now := time.Now()
	ossClient := newMockOSSClient()
	ossClient.SetGetURLs(map[string]string{
		"avatar-1": "https://cdn.example/avatar-1",
		"page-1":   "https://cdn.example/page-1",
		"team-1":   "https://cdn.example/team-1",
	})
	creator := &model.UserInfo{ID: "user-1", Name: "Tester", QQ: "100001", AvatarKey: "avatar-1", IsAvatarUploaded: true, LastLoginAt: now, CreatedAt: now, UpdatedAt: now}
	workset := &model.WorksetInfo{ID: "workset-1", TeamID: "team-1", Name: "WS", Description: "desc", Index: 1, CreatedAt: now, UpdatedAt: now}
	comic := &model.ComicInfo{ID: "comic-1", WorksetID: "workset-1", Index: 2, Title: "Comic", Author: "A", Description: "D", ChapterCount: 3, CreatorID: "user-1", LastActiveAt: now, CreatedAt: now, UpdatedAt: now, Creator: creator, Workset: workset}
	assembledComic := assembleComicInfo(comic, ossClient)
	if assembledComic.Creator == nil || assembledComic.Creator.AvatarURL != "https://cdn.example/avatar-1" || assembledComic.Workset == nil || assembledComic.Workset.ID != "workset-1" {
		t.Fatalf("unexpected assembled comic: %#v", assembledComic)
	}

	page := &model.PageInfo{ID: "page-id", ChapterID: "chapter-1", Index: 1, OSSKey: "page-1", IsUploaded: true, CreatorID: "user-1", TotalUnitCount: 2, TranslatedUnitCount: 1, ProofreadUnitCount: 1, CreatedAt: now, UpdatedAt: now}
	assembledPage := assemblePageInfo(page, ossClient)
	if assembledPage.ImageURL != "https://cdn.example/page-1" {
		t.Fatalf("unexpected assembled page: %#v", assembledPage)
	}

	uploadedAt := now
	translatingAt := now.Add(time.Second)
	translatedAt := now.Add(2 * time.Second)
	proofreadingAt := now.Add(3 * time.Second)
	proofreadAt := now.Add(4 * time.Second)
	typesettingAt := now.Add(5 * time.Second)
	typesetAt := now.Add(6 * time.Second)
	reviewedAt := now.Add(7 * time.Second)
	publishedAt := now.Add(8 * time.Second)
	chapter := &model.ChapterInfo{ID: "chapter-1", ComicID: "comic-1", IsPinned: true, Index: 1, Subtitle: "Sub", PageCount: 2, TotalUnitCount: 10, TranslatedUnitCount: 8, ProofreadUnitCount: 6, CreatorID: "user-1", CreatedAt: now, UpdatedAt: now, UploadedAt: &uploadedAt, TransalatingAt: &translatingAt, TranslatedAt: &translatedAt, ProofreadingAt: &proofreadingAt, ProofreadAt: &proofreadAt, TypesettingAt: &typesettingAt, TypesetAt: &typesetAt, ReviewedAt: &reviewedAt, PublishedAt: &publishedAt, Creator: creator}
	assembledChapter := assembleChapterInfo(chapter, ossClient)
	if assembledChapter.Creator == nil || assembledChapter.PublishedAt == nil || assembledChapter.UploadedAt == nil {
		t.Fatalf("unexpected assembled chapter: %#v", assembledChapter)
	}

	assembledUser, err := assembleUserInfo(creator, ossClient)
	requireNoErr(t, err)
	if assembledUser.AvatarURL != "https://cdn.example/avatar-1" {
		t.Fatalf("unexpected assembled user: %#v", assembledUser)
	}

	team := &model.TeamInfo{ID: "team-1", Name: "Team", Description: "Desc", AvatarOSSKey: "team-1", IsAvatarUploaded: true, CreatedAt: now, UpdatedAt: now}
	assembledTeam, err := assembleTeamInfo(team, ossClient)
	requireNoErr(t, err)
	if assembledTeam.AvatarURL != "https://cdn.example/team-1" {
		t.Fatalf("unexpected assembled team: %#v", assembledTeam)
	}
}

func TestAssembleUserAndTeamInfoReturnOSSFailures(t *testing.T) {
	now := time.Now()
	ossClient := newMockOSSClient()
	ossClient.SetGetErr(errors.New("boom"))

	if _, err := assembleUserInfo(&model.UserInfo{ID: "user-1", AvatarKey: "avatar", IsAvatarUploaded: true, LastLoginAt: now, CreatedAt: now, UpdatedAt: now}, ossClient); err == nil {
		t.Fatal("expected assembleUserInfo error")
	}
	if _, err := assembleTeamInfo(&model.TeamInfo{ID: "team-1", AvatarOSSKey: "team", IsAvatarUploaded: true, CreatedAt: now, UpdatedAt: now}, ossClient); err == nil {
		t.Fatal("expected assembleTeamInfo error")
	}
}

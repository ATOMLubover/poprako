package adapter

import (
	"labelplus-next-web-be/internal/domain/model"
	"labelplus-next-web-be/internal/domain/repository"
	repository_infra "labelplus-next-web-be/internal/infrastructure/repository"
	"labelplus-next-web-be/internal/infrastructure/repository/query_option"
)

func HandleLoadUserInfo(userRepository repository.UserRepository) model.OnLoadUserInfo {
	return func(userID string) (model.UserInfo, error) {
		return userRepository.Get(nil, query_option.FilterByID(repository_infra.UserTable, userID))
	}
}

func HandleLoadMemberInfo(memberRepository repository.MemberRepository) model.OnLoadMemberInfo {
	return func(teamID string, userID string) (model.MemberInfo, error) {
		return memberRepository.Get(
			nil,
			query_option.MemberQuery().FilterByUserID(userID),
			query_option.MemberQuery().FilterByTeamID(teamID),
		)
	}
}

func HandleLoadComicInfo(comicRepository repository.ComicRepository) model.OnLoadComicInfo {
	return func(comicID string) (model.ComicInfo, error) {
		return comicRepository.Get(
			nil,
			query_option.FilterByID(repository_infra.ComicTable, comicID),
		)
	}
}

func HandleLoadAssignmentInfo(assignmentRepository repository.AssignmentRepository) model.OnLoadAssignmentInfo {
	return func(chapterID string, userID string) (model.AssignmentInfo, error) {
		return assignmentRepository.Get(
			nil,
			query_option.AssignmentQuery().FilterByChapterID(chapterID),
			query_option.AssignmentQuery().FilterByUserID(userID),
		)
	}
}

func HandleLoadChapterInfo(chapterRepository repository.ChapterRepository) model.OnLoadChapterInfo {
	return func(chapterID string) (model.ChapterInfo, error) {
		return chapterRepository.Get(
			nil,
			query_option.FilterByID(repository_infra.ChapterTable, chapterID),
		)
	}
}

func HandleLoadPageInfo(pageRepository repository.PageRepository) model.OnLoadPageInfo {
	return func(pageID string) (model.PageInfo, error) {
		return pageRepository.Get(
			nil,
			query_option.FilterByID(repository_infra.PageTable, pageID),
		)
	}
}

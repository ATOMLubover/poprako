package query_option

import intf "labelplus-next-web-be/internal/domain/repository"

type assignmentQuery struct{}

func AssignmentQuery() assignmentQuery {
	return assignmentQuery{}
}

func (assignmentQuery) FilterByChapterID(chapterID string) intf.QueryOption {
	return func(executor intf.Executor) intf.Executor {
		return executor.Where("assignment_table.chapter_id = ?", chapterID)
	}
}

func (assignmentQuery) FilterByUserID(userID string) intf.QueryOption {
    return func(executor intf.Executor) intf.Executor {
        return executor.Where("assignment_table.user_id = ?", userID)
    }
}

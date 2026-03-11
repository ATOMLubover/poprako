package query_option

import intf "labelplus-next-web-be/internal/domain/repository"

type pageQuery struct{}

func PageQuery() pageQuery {
	return pageQuery{}
}

func (pageQuery) FilterByChapterID(chapterID string) intf.QueryOption {
	return func(executor intf.Executor) intf.Executor {
		return executor.Where("page_table.chapter_id = ?", chapterID)
	}
}

func (pageQuery) OrderByIndexAsc() intf.QueryOption {
	return func(executor intf.Executor) intf.Executor {
		return executor.Order("page_table.index ASC")
	}
}

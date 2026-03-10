package query_option

import intf "labelplus-next-web-be/internal/domain/repository"

type chapterQuery struct{}

func ChapterQuery() chapterQuery {
	return chapterQuery{}
}

// FilterByComicID 精确匹配章节所属漫画 ID
func (chapterQuery) FilterByComicID(comicID string) intf.QueryOption {
	return func(executor intf.Executor) intf.Executor {
		return executor.Where("chapter_table.comic_id = ?", comicID)
	}
}

// OrderByIndexAsc 按章节顺序升序排列
func (chapterQuery) OrderByIndexAsc() intf.QueryOption {
	return func(executor intf.Executor) intf.Executor {
		return executor.Order("chapter_table.index ASC")
	}
}

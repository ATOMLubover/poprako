package query_option

import intf "labelplus-next-web-be/internal/domain/repository"

type unitQuery struct{}

func UnitQuery() unitQuery {
	return unitQuery{}
}

func (unitQuery) FilterByPageID(pageID string) intf.QueryOption {
	return func(executor intf.Executor) intf.Executor {
		return executor.Where("unit_table.page_id = ?", pageID)
	}
}

func (unitQuery) OrderByIndexAsc() intf.QueryOption {
	return func(executor intf.Executor) intf.Executor {
		return executor.Order("unit_table.index ASC")
	}
}

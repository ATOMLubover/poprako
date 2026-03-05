package query_option

import (
	intf "labelplus-next-web-be/internal/domain/repository"
)

func CreatedAtDesc() intf.QueryOption {
	return func(executor intf.Executor) intf.Executor {
		return executor.Order("created_at DESC")
	}
}

func FilterByIDs(ids []string) intf.QueryOption {
	return func(executor intf.Executor) intf.Executor {
		return executor.Where("id IN ?", ids)
	}
}

// Paginate 返回分页查询选项
func Paginate(offset, limit int) intf.QueryOption {
	return func(executor intf.Executor) intf.Executor {
		return executor.Offset(offset).Limit(limit)
	}
}

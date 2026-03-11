package query_option

import (
	"fmt"

	intf "labelplus-next-web-be/internal/domain/repository"
)

// 规则（铁律）：
// 1) FilterByXxx 只能筛选当前主表字段（如 xxx_table.xxx），不能用于跨表字段。
// 2) 跨表筛选必须使用 FilterOnXxx（如 FilterOnUserID 表示 user_table.id）。
// 3) 任何 ID 精确筛选必须复用 common.FilterByID，禁止在独立 query 中重复定义。
// 4) query_option 包内禁止裸字段，所有字段都必须使用完整表名前缀（table.column）。

func FilterByID(table string, id string) intf.QueryOption {
	return func(executor intf.Executor) intf.Executor {
		return executor.Where(fmt.Sprintf("%s.id = ?", table), id)
	}
}

func CreatedAtDesc(table string) intf.QueryOption {
	return func(executor intf.Executor) intf.Executor {
		return executor.Order(fmt.Sprintf("%s.created_at DESC", table))
	}
}

func CreatedAtAsc(table string) intf.QueryOption {
	return func(executor intf.Executor) intf.Executor {
		return executor.Order(fmt.Sprintf("%s.created_at ASC", table))
	}
}

func UpdatedAtDesc(table string) intf.QueryOption {
	return func(executor intf.Executor) intf.Executor {
		return executor.Order(fmt.Sprintf("%s.updated_at DESC", table))
	}
}

func FilterByIDs(table string, ids []string) intf.QueryOption {
	return func(executor intf.Executor) intf.Executor {
		return executor.Where(fmt.Sprintf("%s.id IN ?", table), ids)
	}
}

// Paginate 返回分页查询选项
func Paginate(offset, limit int) intf.QueryOption {
	return func(executor intf.Executor) intf.Executor {
		return executor.Offset(offset).Limit(limit)
	}
}

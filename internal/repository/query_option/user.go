package query_option

import (
	"fmt"

	intf "labelplus-next-web-be/internal/domain/repository"
)

type userQuery struct{}

func UserQuery() userQuery {
	return userQuery{}
}

// FilterByQQ 精确匹配用户 QQ 号
func (userQuery) FilterByQQ(qq string) intf.QueryOption {
	return func(executor intf.Executor) intf.Executor {
		return executor.Where("qq = ?", qq)
	}
}

// FilterByFuzzyName 模糊匹配用户昵称
func (userQuery) FilterByFuzzyName(fuzzyName string) intf.QueryOption {
	return func(executor intf.Executor) intf.Executor {
		pattern := fmt.Sprintf("%%%s%%", fuzzyName) // 在 name 前后添加 % 以实现模糊匹配

		return executor.Where("name LIKE ?", pattern)
	}
}

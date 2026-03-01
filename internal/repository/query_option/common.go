package query_option

import intf "labelplus-next-web-be/internal/domain/repository"

func CreatedAtDesc() []intf.QueryOption {
	return []intf.QueryOption{
		func(executor intf.Executor) intf.Executor {
			return executor.Order("created_at DESC")
		},
	}
}

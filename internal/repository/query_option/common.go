package query_option

import (
	"labelplus-next-web-be/internal/domain/model"
	intf "labelplus-next-web-be/internal/domain/repository"
)

func CreatedAtDesc() []intf.QueryOption {
	return []intf.QueryOption{
		func(executor intf.Executor) intf.Executor {
			return executor.Order("created_at DESC")
		},
	}
}

// FilterByQQ 精确匹配 QQ 号
func FilterByQQ(qq string) intf.QueryOption {
	return func(executor intf.Executor) intf.Executor {
		return executor.Where("qq = ?", qq)
	}
}

// FuzzyFilterByName 模糊匹配用户名（依赖 pg_trgm）
func FuzzyFilterByName(name string) intf.QueryOption {
	return func(executor intf.Executor) intf.Executor {
		return executor.Where("name ILIKE ?", "%"+name+"%")
	}
}

// FilterByUserRole 按角色过滤，匹配对应角色列非空的用户
func FilterByUserRole(role model.RoleFlag) intf.QueryOption {
	columnName := roleToColumnName(role)

	return func(executor intf.Executor) intf.Executor {
		if columnName == "" {
			return executor
		}

		return executor.Where(columnName + " IS NOT NULL")
	}
}

// Paginate 返回分页查询选项
func Paginate(offset, limit int) []intf.QueryOption {
	return []intf.QueryOption{
		func(executor intf.Executor) intf.Executor {
			return executor.Offset(offset)
		},
		func(executor intf.Executor) intf.Executor {
			return executor.Limit(limit)
		},
	}
}

func roleToColumnName(role model.RoleFlag) string {
	switch role {
	case model.RolePictureSource:
		return "assigned_picture_source_at"
	case model.RoleTranslator:
		return "assigned_translator_at"
	case model.RoleProofreader:
		return "assigned_proofreader_at"
	case model.RoleTypesetter:
		return "assigned_typesetter_at"
	case model.RoleReviewer:
		return "assigned_reviewer_at"
	case model.RoleAdmin:
		return "assigned_admin_at"
	case model.RoleSuperAdmin:
		return "assigned_super_admin_at"
	default:
		return ""
	}
}

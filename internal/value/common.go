package value

import "errors"

type PaginationParams struct {
	Offset int `json:"offset" url:"offset"`
	Limit  int `json:"limit" url:"limit"`
}

func (pp *PaginationParams) Validate() error {
	if pp == nil {
		return errors.New("分页参数不能为空")
	}

	if pp.Offset < 0 {
		return errors.New("offset 不能为负数")
	}
	if pp.Limit <= 0 {
		return errors.New("limit 必须大于 0")
	}

	return nil
}

package repo_iface

import "poprako-s/internal/domain/model/aggr"

type SysMailRepo interface {
	Send(cre *aggr.SysMailCre) RepoErr
}

package repo_iface

import (
	"poprako-s/internal/domain/model/aggr"
	"poprako-s/internal/domain/model/enum"
	"poprako-s/internal/domain/model/query"
)

// `AnnouncementRepo` defines persistence contract for team announcements.
type AnnouncementRepo interface {
	// `List` returns announcements by query options.
	// Relation loading is controlled by typed `AnnouncementIncl` variadic arguments.
	List(opt query.ListAnnouncementOpt, inc ...enum.AnnouncementIncl) ([]aggr.Announcement, RepoErr)

	// `Create` inserts one team announcement.
	Create(cre *aggr.AnnouncementCre) (*aggr.Announcement, RepoErr)
}

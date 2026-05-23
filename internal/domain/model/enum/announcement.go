package enum

// `AnnouncementIncl` is the include selector for announcement relation preloads.
type AnnouncementIncl string

const (
	// `AnnouncementInclUser` includes the `user` relation in announcement queries.
	AnnouncementInclUser AnnouncementIncl = "user"
)

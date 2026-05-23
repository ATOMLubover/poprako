package query

// `ListAnnouncementOpt` defines filters and pagination for announcement listing.
type ListAnnouncementOpt struct {
	// `TeamId` is required, as announcements are listed by team scope.
	TeamId string

	Pagi PagiOpt
}

package query

// `ListSysMailOpt` holds list query options for system mail listing.
type ListSysMailOpt struct {
	// `RcvId` is the required receiver id of target system mails.
	RcvId string

	// `UnreadOnly` controls whether only unread rows are returned.
	UnreadOnly bool

	// `Pagi` carries explicit pagination options.
	Pagi PagiOpt
}

package aggr

import "time"

type SysMail struct {
	Id string

	// `RcvId` is the user id of the receiver of this notification.
	RcvId string
	// `Read` is true if the notification has been read by the receiver, and false otherwise.
	Read bool

	Title   string
	Content string

	CreatedAt time.Time
}

type SysMailCre struct {
	Id string

	RcvId string

	Title   string
	Content string
}

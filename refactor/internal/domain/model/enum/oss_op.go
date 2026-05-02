package enum

type OssResTyp string

const (
	OssResUserAvatar OssResTyp = "oss_res:user_avatar"
	OssResTeamAvatar OssResTyp = "oss_res:team_avatar"
	OssResComicCover OssResTyp = "oss_res:comic_cover"
	OssResPageImage  OssResTyp = "oss_res:page_image"
)

type OssOp string

const (
	OssOpCre OssOp = "oss_operation:create"
	OssOpDel OssOp = "oss_operation:delete"
)

type OssMsgStatus string

const (
	OssMsgStatePending    OssMsgStatus = "oss_message_status:pending"
	OssMsgStateProcessing OssMsgStatus = "oss_message_status:processing"
	OssMsgStateCompleted  OssMsgStatus = "oss_message_status:completed"
)

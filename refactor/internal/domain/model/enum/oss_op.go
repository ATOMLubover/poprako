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
	OssMsgStatePend OssMsgStatus = "oss_messsage_status:pending"
	OssMsgStateProc OssMsgStatus = "oss_messsage_status:processing"
	OssMsgStateCmpl OssMsgStatus = "oss_messsage_status:completed"
)

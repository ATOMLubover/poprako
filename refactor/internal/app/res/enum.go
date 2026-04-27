package res

type ErrCode int

const (
	BadRequest  ErrCode = 400
	Forbidden   ErrCode = 403
	ServerError ErrCode = 500
)

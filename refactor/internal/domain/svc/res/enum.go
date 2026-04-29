package svc_res

type ErrCode int

const (
	Unauthorized ErrCode = 401
	BadRequest   ErrCode = 400
	NotFound     ErrCode = 404
	Forbidden    ErrCode = 403
	Conflict     ErrCode = 409
	ServerError  ErrCode = 500
	Unavailable  ErrCode = 503
	Timeout      ErrCode = 504
)

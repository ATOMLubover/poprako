package application

const (
	// 由于大部分错误都应该是属于 Bad Request 类型客户端错误
	// 因此仅单独设置一个内部逻辑导致的错误码
	ErrInternalError = "internal_error"
)

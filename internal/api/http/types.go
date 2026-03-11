package http

type FormatResponse struct {
	Code    int    `json:"code"`
	// 错误信息（如果是错误响应），或者成功提示信息（如果是成功响应）
	Message string `json:"message"`
	// 业务数据（仅在成功响应时返回）
	Data    any    `json:"data,omitempty"`
}

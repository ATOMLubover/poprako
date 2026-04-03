package cfg

type AuthCfg struct {
	// ExpHrs 为 token 的过期时间，单位为小时
	ExpHrs int `json:"exp_hrs"`
	// SecretKey 是用于签名 token 的密钥，必须是一个足够复杂的字符串
	// 不直接从配置文件中读取，而是通过环境变量注入，以增强安全性
	SecretKey string `json:"-"`
}

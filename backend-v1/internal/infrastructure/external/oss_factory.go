package external

import (
	"os"

	intf "labelplus-next-web-be/internal/domain/external"
)

func NewOSSClient() intf.OSSClient {
	platform := os.Getenv("OSS_PLATFORM")
	if platform == "" {
		platform = "r2"
	}

	switch platform {
	case "aliyun":
		return NewAliyunOSSClient()
	case "r2":
		return NewR2OSSClient()
	default:
		panic("无效的 OSS_PLATFORM 环境变量：" + platform)
	}
}

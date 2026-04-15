package oss

import (
	"os"
	"strings"

	iface "poprako-s/internal/domain/ext/oss"
)

func NewClient() iface.Client {
	platform := strings.ToLower(strings.TrimSpace(os.Getenv("OSS_PLATFORM")))
	if platform == "" {
		platform = "r2"
	}

	switch platform {
	case "aliyun":
		return NewAliyunClient()
	case "r2":
		return NewR2Client()
	case "noop":
		return NewNoopClient()
	default:
		panic("无效的 OSS_PLATFORM 环境变量: " + platform)
	}
}

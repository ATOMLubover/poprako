package external

import (
	"os"
	"strings"

	intf "labelplus-next-web-be/internal/domain/external"
)

const (
	ossProviderR2      = "r2"
	ossProviderAliyun  = "aliyun"
	ossProviderEnvName = "OSS_PROVIDER"
)

// NewOSSClient 根据环境变量选择 OSS SDK 实现。
//
// 支持的 OSS_PROVIDER:
// - r2: Cloudflare R2（默认）
// - aliyun: 阿里云 OSS Go SDK v2
func NewOSSClient() intf.OSSClient {
	provider := strings.TrimSpace(strings.ToLower(os.Getenv(ossProviderEnvName)))
	if provider == "" {
		provider = ossProviderR2
	}

	switch provider {
	case ossProviderAliyun:
		return NewAliyunOSSClient()
	case ossProviderR2:
		fallthrough
	default:
		return NewR2OSSClient()
	}
}

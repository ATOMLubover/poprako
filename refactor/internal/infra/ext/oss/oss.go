package oss_infra

import (
	oss_iface "poprako-s/internal/domain/ext/oss"

	"github.com/spf13/viper"
)

func NewOssClient() oss_iface.Client {
	pf := viper.GetString("OSS_PLATFORM")
	if pf == "" {
		panic("[NewOssClient] env OSS_PLATFORM is not set")
	}

	switch OssPlatform(pf) {
	case PlatformR2:
		return NewR2Client()
	default:
		panic("[NewOssClient] unsupported OSS_PLATFORM: " + pf)
	}
}

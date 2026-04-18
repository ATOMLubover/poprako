package oss

import (
	"errors"

	iface "poprako-s/internal/domain/ext/oss"
)

// noopClient 是 OSS 客户端的占位实现，所有操作均返回错误
// 正式上线时须替换为真实实现（如 AWS S3 / 阿里云 OSS）
type noopClient struct{}

// NewNoopClient 返回一个不执行任何操作的 OSS 客户端
func NewNoopClient() iface.Client {
	return &noopClient{}
}

func (c *noopClient) GeneratePutPresignedURL(_ string) (string, error) {
	return "", errors.New("OSS 客户端未实现（noopClient）")
}

func (c *noopClient) GenerateGetPresignedURL(_ string) (string, error) {
	return "", nil
}

func (c *noopClient) Delete(_ string) error {
	return nil
}

func (c *noopClient) DeleteBatch(_ []string) error {
	return nil
}

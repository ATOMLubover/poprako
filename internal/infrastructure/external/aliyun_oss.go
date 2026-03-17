package external

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"

	intf "labelplus-next-web-be/internal/domain/external"
)

// AliyunOSSClient 是阿里云 OSS Go SDK v2 的实现。
type AliyunOSSClient struct {
	client *oss.Client

	bucketName      string
	customDomainCDN string
}

// NewAliyunOSSClient 创建阿里云 OSS 客户端。
//
// 必填环境变量:
// - ALIYUN_OSS_ACCESS_KEY_ID
// - ALIYUN_OSS_ACCESS_KEY_SECRET
// - ALIYUN_OSS_REGION
// - ALIYUN_OSS_BUCKET_NAME
//
// 可选环境变量:
// - ALIYUN_OSS_ENDPOINT（默认根据 region 推导）
//
// 必填环境变量（用于下载直链，不参与签名）:
// - ALIYUN_OSS_CUSTOM_DOMAIN_CDN
func NewAliyunOSSClient() intf.OSSClient {
	accessKeyID := os.Getenv("ALIYUN_OSS_ACCESS_KEY_ID")
	if accessKeyID == "" {
		panic("未设置 ALIYUN_OSS_ACCESS_KEY_ID 环境变量")
	}

	accessKeySecret := os.Getenv("ALIYUN_OSS_ACCESS_KEY_SECRET")
	if accessKeySecret == "" {
		panic("未设置 ALIYUN_OSS_ACCESS_KEY_SECRET 环境变量")
	}

	region := os.Getenv("ALIYUN_OSS_REGION")
	if region == "" {
		panic("未设置 ALIYUN_OSS_REGION 环境变量")
	}

	bucketName := os.Getenv("ALIYUN_OSS_BUCKET_NAME")
	if bucketName == "" {
		panic("未设置 ALIYUN_OSS_BUCKET_NAME 环境变量")
	}

	endpoint := os.Getenv("ALIYUN_OSS_ENDPOINT")
	if endpoint == "" {
		endpoint = fmt.Sprintf("https://oss-%s.aliyuncs.com", region)
	}

	customDomainCDN := strings.TrimSpace(os.Getenv("ALIYUN_OSS_CUSTOM_DOMAIN_CDN"))
	if customDomainCDN == "" {
		panic("未设置 ALIYUN_OSS_CUSTOM_DOMAIN_CDN 环境变量")
	}

	customDomainCDN = strings.TrimRight(customDomainCDN, "/")

	cfg := oss.LoadDefaultConfig().
		WithRegion(region).
		WithEndpoint(endpoint).
		WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKeyID, accessKeySecret))

	return &AliyunOSSClient{
		client:          oss.NewClient(cfg),
		bucketName:      bucketName,
		customDomainCDN: customDomainCDN,
	}
}

// GeneratePutPresignedURL 生成上传对象的预签名 URL。
func (aoc *AliyunOSSClient) GeneratePutPresignedURL(objectKey string, contentType string) (string, error) {
	const exp = 10 * time.Minute

	input := &oss.PutObjectRequest{
		Bucket: oss.Ptr(aoc.bucketName),
		Key:    oss.Ptr(objectKey),
	}

	if strings.TrimSpace(contentType) != "" {
		input.ContentType = oss.Ptr(strings.TrimSpace(contentType))
	} else if detectedContentType := detectImageContentType(objectKey); detectedContentType != "" {
		input.ContentType = oss.Ptr(detectedContentType)
	}

	result, err := aoc.client.Presign(context.TODO(), input, oss.PresignExpires(exp))
	if err != nil {
		return "", fmt.Errorf("生成上传预签名 URL 失败: %w", err)
	}

	return result.URL, nil
}

// GenerateGetPresignedURL 生成读取对象的 URL。
//
// 阿里云下载统一返回 CDN 自定义域名直链，不使用签名 URL。
func (aoc *AliyunOSSClient) GenerateGetPresignedURL(objectKey string) (string, error) {
	normalizedObjectKey := normalizeOSSObjectKey(objectKey)
	if normalizedObjectKey == "" {
		return "", nil
	}

	return fmt.Sprintf("%s/%s", aoc.customDomainCDN, normalizedObjectKey), nil
}

// normalizeOSSObjectKey 兼容历史数据中存储完整 URL 的场景，统一提取对象键。
func normalizeOSSObjectKey(rawObjectKey string) string {
	trimmedObjectKey := strings.TrimSpace(rawObjectKey)
	if trimmedObjectKey == "" {
		return ""
	}

	if parsedURL, err := url.Parse(trimmedObjectKey); err == nil && parsedURL.Scheme != "" && parsedURL.Host != "" {
		return strings.TrimPrefix(parsedURL.EscapedPath(), "/")
	}

	return strings.TrimPrefix(trimmedObjectKey, "/")
}

// Delete 删除单个对象。
func (aoc *AliyunOSSClient) Delete(objectKey string) error {
	const maxRetries = 3
	const retryDelay = 500 * time.Millisecond

	ctx := context.Background()
	var lastErr error

	for attempt := 1; attempt <= maxRetries; attempt++ {
		_, err := aoc.client.DeleteObject(ctx, &oss.DeleteObjectRequest{
			Bucket: oss.Ptr(aoc.bucketName),
			Key:    oss.Ptr(objectKey),
		})
		if err == nil {
			return nil
		}

		lastErr = err

		if attempt < maxRetries {
			time.Sleep(retryDelay)
		}
	}

	return fmt.Errorf("在 %d 次尝试后删除对象失败: %w", maxRetries, lastErr)
}

// DeleteBatch 批量删除对象。
func (aoc *AliyunOSSClient) DeleteBatch(objectKeys []string) error {
	if len(objectKeys) == 0 {
		return nil
	}

	const maxRetries = 3
	const retryDelay = 500 * time.Millisecond

	ctx := context.Background()
	var lastErr error

	objects := make([]oss.DeleteObject, 0, len(objectKeys))
	for _, objectKey := range objectKeys {
		objects = append(objects, oss.DeleteObject{Key: oss.Ptr(objectKey)})
	}

	for attempt := 1; attempt <= maxRetries; attempt++ {
		_, err := aoc.client.DeleteMultipleObjects(ctx, &oss.DeleteMultipleObjectsRequest{
			Bucket: oss.Ptr(aoc.bucketName),
			Delete: &oss.Delete{
				Objects: objects,
				Quiet:   true,
			},
		})
		if err == nil {
			return nil
		}

		lastErr = err

		if attempt < maxRetries {
			time.Sleep(retryDelay)
		}
	}

	return fmt.Errorf("在 %d 次尝试后批量删除对象失败: %w", maxRetries, lastErr)
}

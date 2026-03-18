package external

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"

	intf "labelplus-next-web-be/internal/domain/external"
)

type aliyunOSSClient struct {
	client *oss.Client

	bucketName   string
	customDomain string
}

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

	customDomain := strings.TrimSpace(os.Getenv("ALIYUN_OSS_CUSTOM_DOMAIN"))
	if customDomain == "" {
		panic("未设置 ALIYUN_OSS_CUSTOM_DOMAIN 环境变量")
	}

	customDomain = strings.TrimSuffix(strings.TrimPrefix(customDomain, "https://"), "/")
	customDomain = strings.TrimSuffix(strings.TrimPrefix(customDomain, "http://"), "/")

	cfg := oss.LoadDefaultConfig().
		WithRegion(region).
		WithEndpoint(endpoint).
		WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKeyID, accessKeySecret))

	return &aliyunOSSClient{
		client:       oss.NewClient(cfg),
		bucketName:   bucketName,
		customDomain: customDomain,
	}
}

func (aoc *aliyunOSSClient) GeneratePutPresignedURL(objectKey string) (string, error) {
	const exp = 10 * time.Minute

	input := &oss.PutObjectRequest{
		Bucket: oss.Ptr(aoc.bucketName),
		Key:    oss.Ptr(objectKey),
	}

	if contentType := detectImageContentType(objectKey); contentType != "" {
		input.ContentType = oss.Ptr(contentType)
	}

	result, err := aoc.client.Presign(context.TODO(), input, oss.PresignExpires(exp))
	if err != nil {
		return "", fmt.Errorf("生成上传预签名 URL 失败: %w", err)
	}

	return result.URL, nil
}

func (aoc *aliyunOSSClient) GenerateGetPresignedURL(objectKey string) (string, error) {
	if aoc.customDomain != "" {
		return fmt.Sprintf("https://%s/%s", aoc.customDomain, objectKey), nil
	}

	return "", fmt.Errorf("未配置自定义域名")
}

func (aoc *aliyunOSSClient) Delete(objectKey string) error {
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

func (aoc *aliyunOSSClient) DeleteBatch(objectKeys []string) error {
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

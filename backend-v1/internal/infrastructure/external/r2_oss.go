package external

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"

	intf "labelplus-next-web-be/internal/domain/external"
)

type r2OSSClient struct {
	client        *s3.Client
	presignClient *s3.PresignClient

	bucketName   string
	customDomain string
}

func NewR2OSSClient() intf.OSSClient {
	accountID := os.Getenv("R2_ACCOUNT_ID")
	if accountID == "" {
		panic("未设置 R2_ACCOUNT_ID 环境变量")
	}

	accessKeyID := os.Getenv("R2_ACCESS_KEY_ID")
	if accessKeyID == "" {
		panic("未设置 R2_ACCESS_KEY_ID 环境变量")
	}

	secretKeyID := os.Getenv("R2_SECRET_ACCESS_KEY")
	if secretKeyID == "" {
		panic("未设置 R2_SECRET_ACCESS_KEY 环境变量")
	}

	region := os.Getenv("R2_REGION")
	if region == "" {
		region = "auto"
	}

	bucketName := os.Getenv("R2_BUCKET_NAME")
	if bucketName == "" {
		panic("未设置 R2_BUCKET_NAME 环境变量")
	}

	customDomain := os.Getenv("R2_CUSTOM_DOMAIN")
	r2Endpoint := fmt.Sprintf("https://%s.r2.cloudflarestorage.com", accountID)

	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion(region),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(accessKeyID, secretKeyID, ""),
		),
	)
	if err != nil {
		panic(fmt.Sprintf("加载 SDK 配置失败: %v", err))
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(r2Endpoint)
	})

	return &r2OSSClient{
		client:        client,
		presignClient: s3.NewPresignClient(client),
		bucketName:    bucketName,
		customDomain:  customDomain,
	}
}

func (r2 *r2OSSClient) GeneratePutPresignedURL(objectKey string) (string, error) {
	const exp = 10 * time.Minute

	input := &s3.PutObjectInput{
		Bucket: aws.String(r2.bucketName),
		Key:    aws.String(objectKey),
	}

	if contentType := detectImageContentType(objectKey); contentType != "" {
		input.ContentType = aws.String(contentType)
	}

	req, err := r2.presignClient.PresignPutObject(context.TODO(), input, s3.WithPresignExpires(exp))
	if err != nil {
		return "", fmt.Errorf("生成上传预签名 URL 失败: %w", err)
	}

	return req.URL, nil
}

func (r2 *r2OSSClient) GenerateGetPresignedURL(objectKey string) (string, error) {
	if r2.customDomain != "" {
		return fmt.Sprintf("https://%s/%s", r2.customDomain, objectKey), nil
	}

	return "", fmt.Errorf("未配置自定义域名")
}

func (r2 *r2OSSClient) Delete(objectKey string) error {
	const maxRetries = 3
	const retryDelay = 500 * time.Millisecond

	ctx := context.Background()

	var lastErr error

	for attempt := 1; attempt <= maxRetries; attempt++ {
		_, err := r2.client.DeleteObject(ctx, &s3.DeleteObjectInput{
			Bucket: aws.String(r2.bucketName),
			Key:    aws.String(objectKey),
		})
		if err == nil {
			return nil
		}

		var noSuchKey *types.NoSuchKey

		if errors.As(err, &noSuchKey) {
			return nil
		}

		lastErr = err

		if attempt < maxRetries {
			time.Sleep(retryDelay)
		}
	}

	return fmt.Errorf("在 %d 次尝试后删除对象失败: %w", maxRetries, lastErr)
}

func (r2 *r2OSSClient) DeleteBatch(objectKeys []string) error {
	if len(objectKeys) == 0 {
		return nil
	}
	const maxRetries = 3
	const retryDelay = 500 * time.Millisecond

	ctx := context.Background()
	var lastErr error

	objects := make([]types.ObjectIdentifier, 0, len(objectKeys))
	for _, k := range objectKeys {
		objects = append(objects, types.ObjectIdentifier{Key: aws.String(k)})
	}

	for attempt := 1; attempt <= maxRetries; attempt++ {
		out, err := r2.client.DeleteObjects(ctx, &s3.DeleteObjectsInput{
			Bucket: aws.String(r2.bucketName),
			Delete: &types.Delete{
				Objects: objects,
				Quiet:   aws.Bool(true),
			},
		})

		if err == nil {
			if len(out.Errors) == 0 {
				return nil
			}

			var nonNotFound []types.Error
			for _, e := range out.Errors {
				if e.Code != nil && *e.Code == "NoSuchKey" {
					continue
				}
				nonNotFound = append(nonNotFound, e)
			}

			if len(nonNotFound) == 0 {
				return nil
			}

			lastErr = fmt.Errorf("部分对象删除失败: %v", nonNotFound)
		} else {
			var noSuchKey *types.NoSuchKey
			if errors.As(err, &noSuchKey) {
				return nil
			}
			lastErr = err
		}

		if attempt < maxRetries {
			time.Sleep(retryDelay)
		}
	}

	return fmt.Errorf("在 %d 次尝试后批量删除对象失败: %w", maxRetries, lastErr)
}

func detectImageContentType(key string) string {
	ext := strings.ToLower(filepath.Ext(key))
	if ext == ".jpg" || ext == ".jpeg" {
		return "image/jpeg"
	}

	switch ext {
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	case ".webp":
		return "image/webp"
	case ".svg":
		return "image/svg+xml"
	case ".avif":
		return "image/avif"
	case ".bmp":
		return "image/bmp"
	case ".tif", ".tiff":
		return "image/tiff"
	default:
		return ""
	}
}

package oss_infra

import (
	"context"
	"errors"
	"fmt"
	"time"

	oss_iface "poprako-s/internal/domain/ext/oss"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/spf13/viper"
)

type r2Client struct {
	cli    *s3.Client
	preCli *s3.PresignClient

	bucket string
	domain string
}

func NewR2Client() oss_iface.Client {
	accId := viper.GetString("R2_ACCOUNT_ID")
	if accId == "" {
		panic("[NewR2Client] env R2_ACCOUNT_ID is not set")
	}

	acsKeyId := viper.GetString("R2_ACCESS_KEY_ID")
	if acsKeyId == "" {
		panic("[NewR2Client] env R2_ACCESS_KEY_ID is not set")
	}

	scrKeyId := viper.GetString("R2_SECRET_ACCESS_KEY")
	if scrKeyId == "" {
		panic("[NewR2Client] env R2_SECRET_ACCESS_KEY is not set")
	}

	region := viper.GetString("R2_REGION")
	if region == "" {
		region = "auto"
	}

	bucket := viper.GetString("R2_BUCKET_NAME")
	if bucket == "" {
		panic("[NewR2Client] env R2_BUCKET_NAME is not set")
	}

	dom := viper.GetString("R2_CUSTOM_DOMAIN")
	ep := fmt.Sprintf("https://%s.r2.cloudflarestorage.com", accId)

	cfg, err := config.LoadDefaultConfig(
		context.Background(),
		config.WithRegion(region),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(acsKeyId, scrKeyId, ""),
		),
	)
	if err != nil {
		panic(fmt.Sprintf("加载 SDK 配置失败: %v", err))
	}

	cli := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(ep)
	})

	return &r2Client{
		cli:    cli,
		preCli: s3.NewPresignClient(cli),
		bucket: bucket,
		domain: dom,
	}
}

func (c *r2Client) GenGetUrl(key string) (string, error) {
	if c.domain == "" {
		return "", errors.New("[r2Client.GenGetUrl] Non custom domain implementation is not supported")
	}

	return fmt.Sprintf("%s/%s", c.domain, key), nil
}

func (c *r2Client) GenPutUrl(key string) (string, error) {
	const EXP = 10 * time.Minute

	typ := detectImgContTyp(key)
	if typ == "" {
		return "", fmt.Errorf("[r2Client.GenPutUrl] unsupported file type for key: %s", key)
	}

	in := &s3.PutObjectInput{
		Bucket:      aws.String(c.bucket),
		Key:         aws.String(key),
		ContentType: aws.String(typ),
	}

	req, err := c.preCli.PresignPutObject(context.Background(), in, s3.WithPresignExpires(EXP))
	if err != nil {
		return "", fmt.Errorf("[r2Client.GenPutUrl] failed to generate presigned put url: %v", err)
	}

	return req.URL, nil
}

func (c *r2Client) DelBatch(keys []string) error {
	// At current stage, we do not use a exponential backoff strategy for retrying failed deletions, as the number of keys in a batch is expected to be small (usually less than 10), and the likelihood of transient errors is relatively low. However, if we encounter a failure in deleting a batch of keys, we will log the error and the keys that failed to be deleted for further investigation. If we find that transient errors are common in our use case,
	// we can consider implementing a retry mechanism with exponential backoff in the future.
	const MAX_RETRY = 3
	const RETRY_DELAY = time.Second

	objs := make([]types.ObjectIdentifier, len(keys))
	for i, key := range keys {
		objs[i] = types.ObjectIdentifier{Key: aws.String(key)}
	}

	var lastErr error

	for att := range MAX_RETRY {
		out, err := c.cli.DeleteObjects(context.Background(), &s3.DeleteObjectsInput{
			Bucket: aws.String(c.bucket),
			Delete: &types.Delete{
				Objects: objs,
				Quiet:   aws.Bool(true),
			},
		})

		if err != nil && errors.As(err, NoSuchKey()) {
			// Keep silent if a object is already deleted or does not exist,
			// as the end state is the same (the object is not present in the bucket).
			return nil
		}

		if err != nil {
			lastErr = err
			maySleep(att, MAX_RETRY, RETRY_DELAY)

			continue
		}

		// Branch: err == nil.
		if len(out.Errors) == 0 {
			return nil
		}

		// Filter out non-notfound errrors.
		nonNotFound := make([]types.Error, 0)
		for _, e := range out.Errors {
			if e.Code != nil && *e.Code != "NoSuchKey" {
				nonNotFound = append(nonNotFound, e)
			}
		}

		if len(nonNotFound) == 0 {
			return nil
		}

		lastErr = fmt.Errorf("[r2Client.DelBatch] delete batch partially failed: %v", nonNotFound)

		maySleep(att, MAX_RETRY, RETRY_DELAY)
	}

	return lastErr
}

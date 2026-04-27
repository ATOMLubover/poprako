package oss_infra

import (
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

func detectImgContTyp(key string) string {
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

// `NoSchKey` returns a non-nil pointer to `types.NoSuchKey` error,
// which is used to indicate that a specified key does not exist in the S3 bucket. This function can be used in scenarios where you want to simulate or handle the case of a missing NoSuchKey
// without actually performing an S3 operation that would trigger this error.
func NoSuchKey() *types.NoSuchKey {
	return &types.NoSuchKey{}
}

func maySleep(att, max int, dur time.Duration) {
	if att <= 0 {
		return
	}

	if att >= max {
		return
	}

	time.Sleep(dur)
}

package oss_infra

import (
	"errors"
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

func IsNoSuchKey(err error) bool {
	var nsk *types.NoSuchKey
	return err != nil && errors.As(err, &nsk)
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

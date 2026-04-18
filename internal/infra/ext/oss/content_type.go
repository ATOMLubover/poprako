package oss

import (
	"path/filepath"
	"strings"
)

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

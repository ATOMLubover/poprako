package util

import "time"

func NowMillis() int64 {
	return time.Now().UnixMilli()
}

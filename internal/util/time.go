package util

import "time"

func NowMillis() int64 {
	return time.Now().UnixMilli()
}

func ToUnixPtr(t *time.Time) *int64 {
	if t == nil {
		return nil
	}

	unix := t.Unix()

	return &unix
}

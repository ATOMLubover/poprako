package service

import (
	"fmt"
	"math/rand/v2"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// GenID 提供一个不会返回错误的 ID string
func GenID(prefix string) string {
	u, err := uuid.NewV7()
	if err != nil {
		// 回退到简单 rand 生成
		zap.L().Warn(
			"[GenID] UUID 生成失败，回退到简单随机生成",
			zap.Error(err),
		)

		return fmt.Sprintf(
			"%s-%d-%d",
			prefix,
			time.Now().UnixNano(),
			rand.Int32(),
		)
	}

	return fmt.Sprintf("%s-%s", prefix, u.String())
}

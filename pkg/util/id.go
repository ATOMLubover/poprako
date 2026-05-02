package util

import (
	"fmt"
	"math/rand/v2"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// `GenId` always returns a id with no error.
// NOTE: id conflict is possible.
func GenId(prefix string) string {
	u, err := uuid.NewV7()
	if err != nil {
		// 回退到简单 rand 生成
		zap.L().Warn(
			"[GenID] failed to generate uuid, fallback to rand",
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

package lgr

import (
	"os"

	"poprako-s/internal/cfg"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Init(config *cfg.AppCfg) {
	var core zapcore.Core

	if config.IsDevelopment() {
		applyDevelopmentCore(&core)
	} else {
		applyProductionCore(&core)
	}

	// 构建 logger 实例，并注入全局
	logger := zap.New(
		core,
		zap.AddCaller(),                   // 添加行号和文件名信息
		zap.AddStacktrace(zap.PanicLevel), // 在 panic 级别时添加堆栈信息
	)

	zap.ReplaceGlobals(logger)
}

func applyDevelopmentCore(core *zapcore.Core) {
	// 启用默认的 dev 日志配置
	config := zap.NewDevelopmentEncoderConfig()

	// 开启颜色输出
	config.EncodeLevel = zapcore.CapitalColorLevelEncoder
	// 调整时间格式为本地时间
	config.EncodeTime = zapcore.ISO8601TimeEncoder

	// 输出到控制台
	cnslEnc := zapcore.NewConsoleEncoder(config)

	// 默认使用 debug 级别日志
	*core = zapcore.NewCore(
		cnslEnc,
		zapcore.AddSync(os.Stdout),
		zap.DebugLevel,
	)
}

func applyProductionCore(core *zapcore.Core) {
	// 生产环境使用默认的 prod 日志配置
	// 启用 JSON 输出、日志轮转
	config := zap.NewProductionEncoderConfig()

	// 启用方便工具库处理的时间格式
	config.EncodeTime = zapcore.ISO8601TimeEncoder

	// 以 JSON 格式输出
	jsonEnc := zapcore.NewJSONEncoder(config)

	// 配置日志轮转
	lumberLgr := &lumberjack.Logger{
		Filename:   "logs/main-service.log",
		MaxSize:    50, // MB
		MaxBackups: 3,
		MaxAge:     7, // days
		Compress:   true,
	}

	// 同时输出到控制台和日志文件
	syncer := zapcore.NewMultiWriteSyncer(
		zapcore.AddSync(os.Stdout),
		zapcore.AddSync(lumberLgr),
	)

	*core = zapcore.NewCore(
		jsonEnc,
		syncer,
		zap.WarnLevel,
	)
}

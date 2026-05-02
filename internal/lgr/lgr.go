package lgr

import (
	"os"

	"poprako-s/internal/cfg"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func SetGlobal(lgr *zap.Logger) {
	zap.ReplaceGlobals(lgr)
}

func New(appCfg *cfg.AppCfg) *zap.Logger {
	var core zapcore.Core

	switch appCfg.Env {
	case cfg.EnvDev:
		core = newDevCore()
	case cfg.EnvProd:
		core = newProdCore()
	default:
		// Default to development core if environment is unrecognized.
		core = newDevCore()
	}

	lgr := zap.New(
		core,
		zap.AddCaller(),
		zap.AddStacktrace(zap.PanicLevel),
	)

	return lgr
}

func newDevCore() zapcore.Core {
	cfg := zap.NewDevelopmentEncoderConfig()

	cfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
	cfg.EncodeTime = zapcore.ISO8601TimeEncoder

	enc := zapcore.NewConsoleEncoder(cfg)

	core := zapcore.NewCore(
		enc,
		zapcore.AddSync(os.Stdout),
		zap.DebugLevel,
	)

	return core
}

func newProdCore() zapcore.Core {
	cfg := zap.NewProductionEncoderConfig()

	cfg.EncodeTime = zapcore.ISO8601TimeEncoder

	enc := zapcore.NewJSONEncoder(cfg)

	lumberLgr := &lumberjack.Logger{
		Filename:   "logs/main-service.log",
		MaxSize:    50, // MB
		MaxBackups: 3,
		MaxAge:     7, // days
		Compress:   true,
	}

	syncer := zapcore.NewMultiWriteSyncer(
		zapcore.AddSync(os.Stdout),
		zapcore.AddSync(lumberLgr),
	)

	core := zapcore.NewCore(
		enc,
		syncer,
		zap.WarnLevel,
	)

	return core
}

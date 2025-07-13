package log

import (
	"context"
	"os"

	"github.com/yiran15/api-server/base/conf"
	"github.com/yiran15/api-server/base/helper"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger() {
	var (
		encoder  zapcore.Encoder
		writer   zapcore.WriteSyncer
		logLevel zapcore.Level
	)
	logLevelStr := conf.GetServerLogLevel()
	config := zapcore.EncoderConfig{
		TimeKey:       "time",
		LevelKey:      "level",
		NameKey:       "logger",
		CallerKey:     "caller",
		MessageKey:    "msg",
		StacktraceKey: "stacktrace",
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel:   zapcore.CapitalLevelEncoder,
		EncodeTime:    zapcore.ISO8601TimeEncoder,
		EncodeCaller:  zapcore.ShortCallerEncoder,
	}
	encoder = zapcore.NewJSONEncoder(config)
	writer = zapcore.AddSync(os.Stderr)

	switch logLevelStr {
	case "debug":
		logLevel = zap.DebugLevel
	case "info":
		logLevel = zap.InfoLevel
	default:
		logLevel = zap.InfoLevel
	}
	core := zapcore.NewCore(encoder, writer, logLevel)
	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.FatalLevel))
	zap.ReplaceGlobals(logger)
	if logLevel == zap.DebugLevel {
		zap.L().Debug("log initialization successful", zap.String("level", logLevelStr))
	}
}

func WithRequestID(ctx context.Context) *zap.Logger {
	return zap.L().With(zap.String("request_id", helper.GetRequestIDFromContext(ctx)))
}

func WithBody(ctx context.Context, body any) *zap.Logger {
	return WithRequestID(ctx).With(zap.Any("body", body))
}

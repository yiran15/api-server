package log

import (
	"context"
	"os"

	"github.com/yiran15/api-server/base/conf"
	"github.com/yiran15/api-server/base/constant"

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
	writer = zapcore.AddSync(os.Stdout)

	switch logLevelStr {
	case "debug":
		logLevel = zap.DebugLevel
	case "info":
		logLevel = zap.InfoLevel
	default:
		logLevel = zap.InfoLevel
	}
	core := zapcore.NewCore(encoder, writer, logLevel)
	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.FatalLevel), zap.ErrorOutput(os.Stderr))
	zap.ReplaceGlobals(logger)
	if logLevel == zap.DebugLevel {
		zap.L().Debug("log initialization successful", zap.String("level", logLevelStr))
	}
}

func GetRequestID(c context.Context) string {
	if requestID := c.Value(constant.RequestID); requestID != nil {
		return requestID.(string)
	}
	return ""
}

func LWithRequestID(ctx context.Context) *zap.Logger {
	return zap.L().With(zap.String(constant.RequestID, GetRequestID(ctx)))
}

func LWithBody(ctx context.Context, body any) *zap.Logger {
	return LWithRequestID(ctx).With(zap.Any("body", body))
}

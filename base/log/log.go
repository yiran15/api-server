package log

import (
	"context"
	golog "log"
	"os"
	"time"

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
	logLevelStr := conf.GetLogLevel()
	timeZone := conf.GetServerTimeZone()
	cst, err := time.LoadLocation(timeZone)
	if err != nil {
		golog.Printf("failed to load location %s: %v, use local time instead", timeZone, err)
		cst = time.Local
	}

	config := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		EncodeDuration: zapcore.SecondsDurationEncoder,
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.In(cst).Format(time.RFC3339))
		},
		EncodeCaller: zapcore.ShortCallerEncoder,
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
	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.FatalLevel))
	zap.ReplaceGlobals(logger)
	zap.L().Info("log initialization successful", zap.String("level", logLevelStr))
}

func WithRequestID(ctx context.Context) *zap.Logger {
	return zap.L().With(zap.String("request-id", helper.GetRequestIDFromContext(ctx)))
}

func WithBody(ctx context.Context, body any) *zap.Logger {
	return WithRequestID(ctx).With(zap.Any("body", body))
}

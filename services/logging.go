package services

import (
	"github.com/TheZeroSlave/zapsentry"
	"github.com/getsentry/sentry-go"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger(log *zap.Logger, sentryClient *sentry.Client) *zap.Logger {
	cfg := zapsentry.Configuration{
		Level:             zapcore.ErrorLevel, //when to send message to sentry
		EnableBreadcrumbs: true,               // enable sending breadcrumbs to Sentry
		BreadcrumbLevel:   zapcore.InfoLevel,  // at what level should we sent breadcrumbs to sentry
		Tags: map[string]string{
			"service": "api-service",
		},
	}

	core, err := zapsentry.NewCore(cfg, zapsentry.NewSentryClientFromClient(sentryClient))

	//in case of err it will return noop core. so we can safely attach it
	if err != nil {
		log.Fatal("failed to init zap", zap.Error(err))
	}

	zapLog := zapsentry.AttachCoreToLogger(core, log)

	// to use breadcrumbs feature - create new scope explicitly
	// and attach after attaching the core
	return zapLog.With(zapsentry.NewScope())
}

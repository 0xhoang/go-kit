package must

import (
	"github.com/TheZeroSlave/zapsentry"
	"github.com/getsentry/sentry-go"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger(sentryDSN string, serviceName string) (*zap.Logger, *sentry.Client, error) {
	zapLog, _ := zap.NewProduction()

	sentryClient, err := sentry.NewClient(sentry.ClientOptions{
		Dsn: sentryDSN,
	})

	if err != nil {
		return nil, nil, err
	}

	// Flush buffered events before the program terminates.
	// Set the timeout to the maximum duration the program can afford to wait.
	cfg := zapsentry.Configuration{
		Level:             zapcore.ErrorLevel, //when to send message to sentry
		EnableBreadcrumbs: true,               // enable sending breadcrumbs to Sentry
		BreadcrumbLevel:   zapcore.InfoLevel,  // at what level should we sent breadcrumbs to sentry
		Tags: map[string]string{
			"service": serviceName,
		},
	}

	core, err := zapsentry.NewCore(
		cfg,
		zapsentry.NewSentryClientFromClient(sentryClient),
	)

	if err != nil {
		return nil, nil, err
	}

	zapAttach := zapsentry.AttachCoreToLogger(core, zapLog)
	return zapAttach.With(zapsentry.NewScope()), sentryClient, nil
}

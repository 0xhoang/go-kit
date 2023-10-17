package main

import (
	"fmt"
	"github.com/0xhoang/go-kit/cmd/data"
	"github.com/0xhoang/go-kit/common"
	"github.com/0xhoang/go-kit/config"
	"github.com/0xhoang/go-kit/services"
	"github.com/getsentry/sentry-go"
	"go.uber.org/zap"
	"time"

	log "github.com/sirupsen/logrus"
	"os"
)

func main() {
	cfg := config.ReadConfigAndArg()

	sentryClient, err := sentry.NewClient(sentry.ClientOptions{
		Dsn: cfg.SentryDSN,
	})

	if err != nil {
		log.Fatal("failed to init sentry", zap.Error(err))
	}
	// Flush buffered events before the program terminates.
	// Set the timeout to the maximum duration the program can afford to wait.
	defer sentryClient.Flush(2 * time.Second)

	zapLog, _ := zap.NewProduction()
	defer zapLog.Sync()
	logger := services.NewLogger(zapLog, sentryClient)

	//db := database.ConnectDb(cfg.Db)

	helloSvc := services.NewHelloService(
		logger,
		cfg,
		nil,
	)

	rootCmd := data.NewRootCmd(helloSvc)

	if err := rootCmd.Execute(); err != nil {
		fmt.Print(common.Red)
		fmt.Println("Execute", err.Error())
		fmt.Print(common.Reset)
		os.Exit(1)
	}
}

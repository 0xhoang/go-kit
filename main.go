package main

import (
	"context"
	"fmt"
	"github.com/0xhoang/go-kit/api"
	"github.com/0xhoang/go-kit/common"
	"github.com/0xhoang/go-kit/config"
	"github.com/0xhoang/go-kit/dao/payment"
	"github.com/0xhoang/go-kit/dao/users"
	"github.com/0xhoang/go-kit/database"
	"github.com/0xhoang/go-kit/services"
	"github.com/0xhoang/go-kit/task"
	"github.com/allegro/bigcache/v3"
	"github.com/getsentry/sentry-go"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	common.SwaggerConfig()

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

	db := database.ConnectDb(cfg.Db)
	err = database.Migration(db)
	if err != nil {
		log.Fatalf("migration: %v", err)
	}

	if err := database.AutoSeedingData(db); err != nil {
		//log.Fatalf("seeding: %v", err)
	}

	life := time.Hour * 24 * 7 * 52 * 10 //10year
	_, _ = bigcache.NewBigCache(bigcache.DefaultConfig(life))

	//start server
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"GET", "POST", "PUT", "HEAD", "OPTIONS", "DELETE"},
		AllowHeaders:    []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		MaxAge:          12 * time.Hour,
	}))

	userDao := users.NewUser(db)
	paymentDao := payment.NewPaymentAddressAction(db)

	helloServiceSvc := services.NewHelloService(logger, cfg, db)
	userSvc := services.NewUser(logger, cfg, userDao)

	var cronJob = cron.New(cron.WithParser(cron.NewParser(
		cron.SecondOptional | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor,
	)))

	eventSvc := task.NewEventService(logger, cronJob, db, cfg, paymentDao)

	svr := api.NewServer(
		r,
		nil,
		logger,
		helloServiceSvc,
		userSvc,
		eventSvc,
	)

	authMw := svr.AuthMiddleware("key-secret")
	svr.WithAuthMw(authMw)
	svr.Routes()

	//todo: worker
	go eventSvc.StartEventPaymentAction()
	cronJob.Start()

	port := cfg.Port
	address := fmt.Sprintf(":%d", port)
	srv := &http.Server{
		Addr:    address,
		Handler: r,
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		fmt.Printf("Listening and serving HTTP on %s\n", address)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("router.Run: %s\n", zap.Error(err))
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Error("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown:", zap.Error(err))
	}

	logger.Info("Server exiting")
}

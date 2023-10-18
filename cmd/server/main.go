package main

import (
	"context"
	"fmt"
	"github.com/0xhoang/go-kit/cmd/task"
	"github.com/0xhoang/go-kit/common"
	"github.com/0xhoang/go-kit/config"
	"github.com/0xhoang/go-kit/internal/dao"
	"github.com/0xhoang/go-kit/internal/must"
	"github.com/0xhoang/go-kit/internal/services"
	"github.com/0xhoang/go-kit/migration"
	"github.com/allegro/bigcache/v3"
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

var cronJob = cron.New(cron.WithParser(cron.NewParser(
	cron.SecondOptional | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor,
)))

func initServer(g *gin.Engine, logger *zap.Logger, cfg *config.Config) error {
	db := must.ConnectDb(cfg.Db)
	err := migration.Migration(db)
	if err != nil {
		return err
	}

	if err := migration.AutoSeedingData(db); err != nil {
		//log.Fatalf("seeding: %v", err)
	}

	//dao
	userDao := dao.NewUser(db)
	paymentDao := dao.NewPaymentAddressAction(db)
	//svc
	helloServiceSvc := services.NewHelloService(logger, cfg, db)
	userSvc := services.NewUser(logger, cfg, userDao)
	eventSvc := task.NewEventService(logger, cronJob, db, cfg, paymentDao)

	svr := NewServer(
		g,
		nil,
		logger,
		helloServiceSvc,
		userSvc,
		eventSvc,
	)
	svr.AuthMiddleware("key-secret")
	svr.Routes()

	//todo: worker
	go eventSvc.StartEventPaymentAction()
	cronJob.Start()

	return nil
}

func main() {
	common.SwaggerConfig()

	var ctx = context.TODO()
	cfg := config.ReadConfigAndArg()

	logger, sentry, err := must.NewLogger(cfg.SentryDSN)
	if err != nil {
		log.Fatalf("logger: %v", err)
	}

	defer logger.Sync()
	defer sentry.Flush(2 * time.Second)

	life := time.Hour * 24 * 7 * 52 * 10 //10year
	_, _ = bigcache.NewBigCache(bigcache.DefaultConfig(life))

	gin := gin.Default()
	gin.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"GET", "POST", "PUT", "HEAD", "OPTIONS", "DELETE"},
		AllowHeaders:    []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		MaxAge:          12 * time.Hour,
	}))

	srv := &http.Server{
		Handler: gin,
		Addr:    fmt.Sprintf(":%d", cfg.Port),
	}

	err = initServer(gin, logger, cfg)
	if err != nil {
		log.Fatalf("run: %v", err)
	}

	go func() {
		fmt.Printf("Listening and serving HTTP on %s\n", fmt.Sprintf(":%d", cfg.Port))
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
	ctxTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctxTimeout); err != nil {
		logger.Fatal("Server forced to shutdown:", zap.Error(err))
	}

	logger.Info("Server exiting")
}

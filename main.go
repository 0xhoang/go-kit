package main

import (
	"context"
	"fmt"
	"github.com/allegro/bigcache/v3"
	"github.com/getsentry/sentry-go"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gitlab.com/idolauncher/go-template-kit/api"
	"gitlab.com/idolauncher/go-template-kit/common"
	"gitlab.com/idolauncher/go-template-kit/config"
	"gitlab.com/idolauncher/go-template-kit/dao"
	"gitlab.com/idolauncher/go-template-kit/services"
	"go.uber.org/zap"
	"gorm.io/gorm"
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

	//uncomment if use db
	//db := database.ConnectDb(cfg.MysqlDB)
	//comment if use db
	var db *gorm.DB
	//uncomment if use migration db
	/*err := database.Migration(db)
	if err != nil {
		log.Fatalf("migration: %v", err)
	}*/

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

	_ = dao.NewUser(db)
	helloServiceSvc := services.NewHelloService(logger, cfg, db)

	svr := api.NewServer(
		r,
		nil,
		logger,
		helloServiceSvc)

	svr.AuthMiddleware("key")
	svr.Routes()

	port := cfg.Port
	if cfg.Env != "development" {
		port = 8888
	}

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

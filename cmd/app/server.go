package main

import (
	"context"
	"github.com/0xhoang/go-kit/cmd/task"
	"github.com/0xhoang/go-kit/common"
	"github.com/0xhoang/go-kit/config"
	"github.com/0xhoang/go-kit/internal/dao"
	"github.com/0xhoang/go-kit/internal/must"
	"github.com/0xhoang/go-kit/internal/services"
	"github.com/0xhoang/go-kit/migration"
	"github.com/allegro/bigcache/v3"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	"github.com/robfig/cron/v3"
	"google.golang.org/grpc"
	"log"
	"time"
)

var cronJob = cron.New(cron.WithParser(cron.NewParser(
	cron.SecondOptional | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor,
)))

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

	db := must.ConnectDb(cfg.Db)
	err = migration.Migration(db)
	if err != nil {
		log.Fatalf("migration: %v", err)
	}

	if err := migration.AutoSeedingData(db); err != nil {
		//log.Fatalf("seeding: %v", err)
	}

	_, _ = bigcache.New(context.Background(), bigcache.DefaultConfig(10*time.Minute))

	//dao
	userDao := dao.NewUser(db)
	paymentDao := dao.NewPaymentAddressAction(db)

	eventSvc := task.NewEventService(logger, cronJob, db, cfg, paymentDao)

	go eventSvc.StartEventPaymentAction()
	cronJob.Start()

	middlewareAuth := NewMiddleware(cfg.AuthenticationPubSecretKey)
	opt := []grpc.ServerOption{
		grpc.StreamInterceptor(auth.StreamServerInterceptor(middlewareAuth.AuthMiddleware)),
		grpc.UnaryInterceptor(auth.UnaryServerInterceptor(middlewareAuth.AuthMiddleware)),
	}

	must.NewServer(ctx, cfg,
		opt,
		services.NewGokitService(
			logger,
			cfg,
			db,
			userDao,
		),
		services.NewGokitPublicService(
			logger,
			cfg,
			userDao,
		),
	)
}

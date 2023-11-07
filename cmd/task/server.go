package main

import (
	"github.com/0xhoang/go-kit/config"
	"github.com/0xhoang/go-kit/internal/dao"
	"github.com/0xhoang/go-kit/internal/must"
	"github.com/0xhoang/go-kit/migration"
	"github.com/robfig/cron/v3"
	"log"
	"time"
)

var cronJob = cron.New(cron.WithParser(cron.NewParser(
	cron.SecondOptional | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor,
)))

func main() {
	cfg := config.ReadConfigAndArg()
	logger, sentry, err := must.NewLogger(cfg.SentryDSN, "task")
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

	paymentDao := dao.NewPaymentAddressAction(db)
	eventSvc := NewEventService(logger, cronJob, db, cfg, paymentDao)
	go eventSvc.StartEventPaymentAction()

	cronJob.Start()

	select {}
}

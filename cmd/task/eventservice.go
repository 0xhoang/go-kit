package main

import (
	taskCommon "github.com/0xhoang/go-kit/cmd/task/common"
	"github.com/0xhoang/go-kit/cmd/task/paymentaction"
	"github.com/0xhoang/go-kit/config"
	"github.com/0xhoang/go-kit/internal/dao"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type EventService struct {
	logger *zap.Logger
	job    *cron.Cron
	db     *gorm.DB
	cfg    *config.Config
	dao    dao.PaymentAddressActionDaoInterface
}

func NewEventService(logger *zap.Logger, job *cron.Cron, db *gorm.DB, cfg *config.Config, dao dao.PaymentAddressActionDaoInterface) *EventService {
	return &EventService{logger: logger, job: job, db: db, cfg: cfg, dao: dao}
}

func (s *EventService) StartEventPaymentAction() {
	s.logger.Info("start event")

	historyPool := taskCommon.NewWorkerPool(10, 10, "history")
	//defer historyPool.Shutdown()

	property := &taskCommon.AssetsProperties{
		Logger:      s.logger,
		Db:          s.db,
		Cfg:         s.cfg,
		HistoryPool: historyPool,
		Dao:         s.dao,
	}

	worker := paymentaction.NewPaymentActionTask(property)

	_, err := s.job.AddFunc("@every 5s", func() {
		submittedPool := taskCommon.NewWorkerPool(10, 20, "SubmittedAction")
		defer submittedPool.Shutdown()

		worker.SubmittedAction(submittedPool)
	})

	if err != nil {
		s.logger.Fatal("AddFunc error", zap.Error(err))
	}

	_, err = s.job.AddFunc("@every 5s", func() {
		processingPool := taskCommon.NewWorkerPool(10, 20, "ProcessingAction")
		defer processingPool.Shutdown()

		worker.ProcessingAction(processingPool)
	})

	if err != nil {
		s.logger.Fatal("AddFunc error", zap.Error(err))
	}

	_, err = s.job.AddFunc("@every 5s", func() {
		confirmingPool := taskCommon.NewWorkerPool(10, 20, "ConfirmingAction")
		defer confirmingPool.Shutdown()

		worker.ConfirmingAction(confirmingPool)
	})

	if err != nil {
		s.logger.Fatal("AddFunc error", zap.Error(err))
	}
}

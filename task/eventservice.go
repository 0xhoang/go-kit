package task

import (
	"github.com/0xhoang/go-kit/config"
	"github.com/0xhoang/go-kit/dao/payment"
	taskCommon "github.com/0xhoang/go-kit/task/common"
	"github.com/0xhoang/go-kit/task/paymentaction"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type EventService struct {
	logger *zap.Logger
	job    *cron.Cron
	db     *gorm.DB
	cfg    *config.Config
	dao    *payment.PaymentAddressAction
}

func NewEventService(logger *zap.Logger, job *cron.Cron, db *gorm.DB, cfg *config.Config, dao *payment.PaymentAddressAction) *EventService {
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
		log.Fatalf("AddFunc error = %v", err.Error())
	}

	_, err = s.job.AddFunc("@every 5s", func() {
		processingPool := taskCommon.NewWorkerPool(10, 20, "ProcessingAction")
		defer processingPool.Shutdown()

		worker.ProcessingAction(processingPool)
	})

	if err != nil {
		log.Fatalf("AddFunc error = %v", err.Error())
	}

	_, err = s.job.AddFunc("@every 5s", func() {
		confirmingPool := taskCommon.NewWorkerPool(10, 20, "ConfirmingAction")
		defer confirmingPool.Shutdown()

		worker.ConfirmingAction(confirmingPool)
	})

	if err != nil {
		log.Fatalf("AddFunc error = %v", err.Error())
	}
}

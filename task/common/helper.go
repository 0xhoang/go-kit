package common

import (
	"github.com/0xhoang/go-kit/common/fsm"
	"github.com/0xhoang/go-kit/config"
	"github.com/0xhoang/go-kit/dao/payment"
	"github.com/0xhoang/go-kit/models"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type TaskHelper struct {
}

type AssetsProperties struct {
	Logger      *zap.Logger
	Db          *gorm.DB
	Cfg         *config.Config
	HistoryPool *WorkerPool
	Dao         *payment.PaymentAddressAction
}

// finite-state machine
const (
	SUBMITED_STATE   fsm.StateType = models.AasmStateSubmitted
	PROCESSING_STATE fsm.StateType = models.AasmStateProcessing
	CONFIRMING_STATE fsm.StateType = models.AasmStateConfirming
	SUCCEED_STATE    fsm.StateType = models.AasmStateSucceed
	ERRORED_STATE    fsm.StateType = models.AasmStateErrored

	//Event
	SUBMITED_EVENT   fsm.EventType = "submitted_event"
	PROCESSING_EVENT fsm.EventType = "processing_event"
	SUCCEED_EVENT    fsm.EventType = "succeed_event"
	ERRORED_EVENT    fsm.EventType = "errored_event"
)

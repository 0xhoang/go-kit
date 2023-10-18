package common

import (
	"github.com/0xhoang/go-kit/common"
	"github.com/0xhoang/go-kit/common/fsm"
	"github.com/0xhoang/go-kit/config"
	"github.com/0xhoang/go-kit/internal/dao"
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
	Dao         dao.PaymentAddressActionDaoInterface
}

// finite-state machine
const (
	SUBMITED_STATE   fsm.StateType = common.AasmStateSubmitted
	PROCESSING_STATE fsm.StateType = common.AasmStateProcessing
	CONFIRMING_STATE fsm.StateType = common.AasmStateConfirming
	SUCCEED_STATE    fsm.StateType = common.AasmStateSucceed
	ERRORED_STATE    fsm.StateType = common.AasmStateErrored

	//Event
	SUBMITED_EVENT   fsm.EventType = "submitted_event"
	PROCESSING_EVENT fsm.EventType = "processing_event"
	SUCCEED_EVENT    fsm.EventType = "succeed_event"
	ERRORED_EVENT    fsm.EventType = "errored_event"
)

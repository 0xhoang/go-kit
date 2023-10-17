package paymentaction

import (
	"fmt"
	"github.com/0xhoang/go-kit/common/fsm"
	"github.com/0xhoang/go-kit/models"
	taskCommon "github.com/0xhoang/go-kit/task/common"

	"sync"
)

type ProcessingAction struct {
	taskCommon.TaskHelper
	model    *models.CustodialPaymentAddressAction
	property *taskCommon.AssetsProperties
}

func NewProcessingAction(model *models.CustodialPaymentAddressAction, property *taskCommon.AssetsProperties) *ProcessingAction {
	return &ProcessingAction{model: model, property: property}
}

func (s *ProcessingAction) Task() {
	s.property.Logger.Info("Start run NewProcessingAction...")
	s.property.Logger.Info(fmt.Sprintf("Record id: %d", s.model.ID))

	if s.model == nil {
		s.property.Logger.Info("s.model is nul")
		return
	}

	withdrawFsm := NewAasmState(fsm.StateType(s.model.AasmState))
	if withdrawFsm.Current != taskCommon.PROCESSING_STATE {
		s.TrackHistory(models.LOG_STATUS_FAILED, fmt.Sprintf("aass_state is not true, %v != %v, stage = %v", withdrawFsm.Current, taskCommon.PROCESSING_STATE, s.model.AasmState), "")
		s.property.Dao.UpdateQueueRetry(s.property.Cfg.JobInfo.TimesRetry, s.model)
		return
	}

	processingCtx := &AasmStateContext{
		models: s.model,
		dao:    s.property.Dao,
	}

	s.TrackHistory(models.LOG_STATUS_SUCCESS, fmt.Sprintf("next step = %v", taskCommon.PROCESSING_EVENT), "")

	err := withdrawFsm.SendEvent(taskCommon.PROCESSING_EVENT, processingCtx)
	if err != nil {
		s.TrackHistory(models.LOG_STATUS_FAILED, fmt.Sprintf("SendEvent %v", taskCommon.PROCESSING_EVENT), err.Error())
		s.property.Dao.UpdateQueueRetry(s.property.Cfg.JobInfo.TimesRetry, s.model)
		return
	}
}

func (s *ProcessingAction) TrackHistory(status string, message string, responseData string) {
	trackData := &models.CustodialPaymentAddressLog{
		CustodialPaymentActionID: s.model.ID,
		Action:                   s.model.AasmState,
		State:                    s.model.AasmState,
		Status:                   status,
		Msg:                      message,
		Data:                     responseData,
	}

	np := NewHistoryAction(trackData, s.property)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		// Submit the task to be worked on. When RunTask
		// returns we know it is being handled.
		s.property.HistoryPool.Run(np)
		wg.Done()
	}()
	wg.Wait()
}

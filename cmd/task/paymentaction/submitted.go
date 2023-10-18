package paymentaction

import (
	"fmt"
	taskCommon "github.com/0xhoang/go-kit/cmd/task/common"
	"github.com/0xhoang/go-kit/common"
	"github.com/0xhoang/go-kit/common/fsm"
	"github.com/0xhoang/go-kit/internal/models"
	"sync"
)

type SubmittedAction struct {
	taskCommon.TaskHelper
	model    *models.CustodialPaymentAddressAction
	property *taskCommon.AssetsProperties
}

func NewSubmittedAction(model *models.CustodialPaymentAddressAction, property *taskCommon.AssetsProperties) *SubmittedAction {
	return &SubmittedAction{model: model, property: property}
}

func (s *SubmittedAction) Task() {
	s.property.Logger.Info("Start run NewSubmittedAction...")
	s.property.Logger.Info(fmt.Sprintf("Record id: %d", s.model.ID))

	if s.model == nil {
		s.property.Logger.Info("s.model is nul")
		return
	}

	withdrawFsm := NewAasmState(fsm.StateType(s.model.AasmState))
	if withdrawFsm.Current != taskCommon.SUBMITED_STATE {
		s.TrackHistory(common.LOG_STATUS_FAILED, fmt.Sprintf("aass_state is not true, %v != %v, stage = %v", withdrawFsm.Current, taskCommon.SUBMITED_STATE, s.model.AasmState), "")
		s.property.Dao.UpdateQueueRetry(s.property.Cfg.JobInfo.TimesRetry, s.model)
		return
	}

	// Let's now retry the same order with a valid set of items.
	processingCtx := &AasmStateContext{
		models: s.model,
		dao:    s.property.Dao,
	}

	s.TrackHistory(common.LOG_STATUS_SUCCESS, fmt.Sprintf("next step = %v", taskCommon.SUBMITED_EVENT), "")

	err := withdrawFsm.SendEvent(taskCommon.SUBMITED_EVENT, processingCtx)
	if err != nil {
		s.TrackHistory(common.LOG_STATUS_FAILED, fmt.Sprintf("SendEvent %v", taskCommon.SUBMITED_EVENT), err.Error())
		s.property.Dao.UpdateQueueRetry(s.property.Cfg.JobInfo.TimesRetry, s.model)
		return
	}

	return
}

func (s *SubmittedAction) TrackHistory(status string, message string, responseData string) {
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

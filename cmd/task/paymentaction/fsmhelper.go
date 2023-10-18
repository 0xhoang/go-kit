package paymentaction

import (
	"fmt"
	taskCommon "github.com/0xhoang/go-kit/cmd/task/common"
	"github.com/0xhoang/go-kit/common"
	"github.com/0xhoang/go-kit/common/fsm"
	"github.com/0xhoang/go-kit/internal/dao"
	"github.com/0xhoang/go-kit/internal/models"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type AasmStateContext struct {
	models *models.CustodialPaymentAddressAction
	dao    dao.PaymentAddressActionDaoInterface
	err    error
}

func NewAasmState(current fsm.StateType) *fsm.StateMachine {
	return &fsm.StateMachine{
		Current: current,
		States: fsm.States{
			taskCommon.SUBMITED_STATE: fsm.State{ //current state
				Events: fsm.Events{
					//current event: next state
					taskCommon.SUCCEED_EVENT:  taskCommon.SUCCEED_STATE,
					taskCommon.SUBMITED_EVENT: taskCommon.PROCESSING_STATE,
				},
			},
			taskCommon.PROCESSING_STATE: fsm.State{
				Action: &ProcessEvent{},
				Events: fsm.Events{
					taskCommon.SUCCEED_EVENT:    taskCommon.SUCCEED_STATE,
					taskCommon.PROCESSING_EVENT: taskCommon.CONFIRMING_STATE,
				},
			},
			taskCommon.CONFIRMING_STATE: fsm.State{
				Action: &ConfirmingEvent{},
				Events: fsm.Events{
					taskCommon.SUCCEED_EVENT: taskCommon.SUCCEED_STATE,
					taskCommon.ERRORED_EVENT: taskCommon.ERRORED_STATE,
				},
			},
			taskCommon.SUCCEED_STATE: fsm.State{
				Action: &SucceedEvent{},
			},
			taskCommon.ERRORED_STATE: fsm.State{
				Action: &ErroredEvent{},
			},
		},
	}
}

type ProcessEvent struct{}

func (a *ProcessEvent) Execute(eventCtx fsm.EventContext) (fsm.EventType, error) {
	event := eventCtx.(*AasmStateContext)
	if event.models == nil {
		return "", errors.New("models not found")
	}

	if event.dao == nil {
		return "", errors.New("dao not found")
	}

	log.Info("Execute, id:", event.models.ID)

	withdraw := event.models

	aasmState := fmt.Sprintf("%v", taskCommon.PROCESSING_STATE)
	if err := event.dao.UpdateState(withdraw.ID, aasmState, int(common.StageStatus), 0); err != nil {
		log.Error(fmt.Sprintf("Error updating log: id %v err %v", withdraw.ID, err))
		return "", err
	}

	return fsm.NoOp, nil
}

type ConfirmingEvent struct{}

func (a *ConfirmingEvent) Execute(eventCtx fsm.EventContext) (fsm.EventType, error) {
	event := eventCtx.(*AasmStateContext)
	if event.models == nil {
		return "", errors.New("models not found")
	}

	if event.dao == nil {
		return "", errors.New("dao not found")
	}

	log.Info("Execute, id:", event.models.ID)

	withdraw := event.models
	aasmState := fmt.Sprintf("%v", taskCommon.CONFIRMING_STATE)
	if err := event.dao.UpdateState(withdraw.ID, aasmState, int(common.StageStatus), 0); err != nil {
		log.Error(fmt.Sprintf("Error updating log: id %v err %v", withdraw.ID, err))
		return "", err
	}

	return fsm.NoOp, nil
}

type SucceedEvent struct{}

func (a *SucceedEvent) Execute(eventCtx fsm.EventContext) (fsm.EventType, error) {
	event := eventCtx.(*AasmStateContext)
	if event.models == nil {
		return "", errors.New("models not found")
	}

	if event.dao == nil {
		return "", errors.New("dao not found")
	}

	log.Info("Execute, id:", event.models.ID)

	withdraw := event.models
	aasmState := fmt.Sprintf("%v", taskCommon.SUCCEED_STATE)
	if err := event.dao.UpdateState(withdraw.ID, aasmState, int(common.StageStatus), 0); err != nil {
		log.Error(fmt.Sprintf("Error updating log: id %v err %v", withdraw.ID, err))
		return "", err
	}

	return fsm.NoOp, nil
}

type ErroredEvent struct{}

func (a *ErroredEvent) Execute(eventCtx fsm.EventContext) (fsm.EventType, error) {
	event := eventCtx.(*AasmStateContext)
	if event.models == nil {
		return "", errors.New("models not found")
	}

	if event.dao == nil {
		return "", errors.New("dao not found")
	}

	log.Info("Execute, id:", event.models.ID)

	withdraw := event.models
	aasmState := fmt.Sprintf("%v", taskCommon.ERRORED_STATE)
	if err := event.dao.UpdateState(withdraw.ID, aasmState, int(common.StageStatus), 0); err != nil {
		log.Error(fmt.Sprintf("Error updating log: id %v err %v", withdraw.ID, err))
		return "", err
	}

	return fsm.NoOp, nil
}

package paymentaction

import (
	"fmt"
	"github.com/0xhoang/go-kit/models"
	taskCommon "github.com/0xhoang/go-kit/task/common"
	"go.uber.org/zap"
	"reflect"
	"sync"
)

type PaymentActionTask struct {
	property *taskCommon.AssetsProperties
}

func NewPaymentActionTask(property *taskCommon.AssetsProperties) *PaymentActionTask {
	return &PaymentActionTask{property: property}
}

func (t *PaymentActionTask) addTask(processList []*models.CustodialPaymentAddressAction, workerPool *taskCommon.WorkerPool, taskType taskCommon.Worker) {
	var wg sync.WaitGroup
	wg.Add(len(processList))

	for _, value := range processList {
		t.property.Logger.Info(fmt.Sprintf("Model id: %v\n", value.ID))

		var np taskCommon.Worker
		switch reflect.ValueOf(taskType).Type() {

		case reflect.ValueOf(&SubmittedAction{}).Type():
			np = NewSubmittedAction(
				value,
				t.property,
			)
		case reflect.ValueOf(&ProcessingAction{}).Type():
			np = NewProcessingAction(
				value,
				t.property,
			)
		case reflect.ValueOf(&ConfirmingAction{}).Type():
			np = NewConfirmingAction(
				value,
				t.property,
			)
		default:
			t.property.Logger.Error("unknown task type")
			return
		}

		//update status running
		err := t.property.Dao.UpdateQueueRunning(value)
		if err != nil {
			t.property.Logger.Info(fmt.Sprintf("Received message: %d and update error", err.Error()))
			continue
		}

		go func(item *models.CustodialPaymentAddressAction) {
			// Submit the task to be worked on. When RunTask
			// returns we know it is being handled.
			t.property.Logger.Info(fmt.Sprintf("Received message: %d and push to job", item.ID))
			workerPool.Run(np)
			wg.Done()
		}(value)
	}

	wg.Wait()
}

func (t *PaymentActionTask) SubmittedAction(workerPool *taskCommon.WorkerPool) {
	processList, err := t.property.Dao.ListQueue(
		[]string{
			models.AasmStateSubmitted,
		}, nil,
		models.StageStatus,
	)

	if err != nil {
		t.property.Logger.Error("SubmittedTask", zap.Error(err))
		return
	}

	if len(processList) <= 0 {
		fmt.Println("State SubmittedAction [0] item")
		return
	}

	t.addTask(processList, workerPool, &SubmittedAction{})
}

func (t *PaymentActionTask) ProcessingAction(workerPool *taskCommon.WorkerPool) {
	processList, err := t.property.Dao.ListQueue(
		[]string{
			models.AasmStateProcessing,
		},
		nil,
		models.StageStatus,
	)

	if err != nil {
		t.property.Logger.Error("ProcessingAction", zap.Error(err))
		return
	}

	t.addTask(processList, workerPool, &ProcessingAction{})
}

func (t *PaymentActionTask) ConfirmingAction(workerPool *taskCommon.WorkerPool) {
	processList, err := t.property.Dao.ListQueue(
		[]string{
			models.AasmStateConfirming,
		},
		nil,
		models.StageStatus,
	)

	if err != nil {
		t.property.Logger.Error("ConfirmingAction", zap.Error(err))
		return
	}

	if len(processList) <= 0 {
		fmt.Println("State ConfirmingAction [0] item")
		return
	}

	t.addTask(processList, workerPool, &ConfirmingAction{})
}

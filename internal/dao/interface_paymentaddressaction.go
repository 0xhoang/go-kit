package dao

import (
	"github.com/0xhoang/go-kit/common"
	"github.com/0xhoang/go-kit/internal/models"
)

type PaymentAddressActionDaoInterface interface {
	Update(model *models.CustodialPaymentAddressAction) error
	GetJobProcessing(userId uint64) ([]*models.CustodialPaymentAddressAction, error)
	CreateJob(model *models.CustodialPaymentAddressAction) error
	UpdateRawData(id uint, rawData, txHash string) error
	UpdateState(id uint, aasmState string, stageStatus int, errCount int) error
	InsertLog(data *models.CustodialPaymentAddressLog) error
	ListQueue(aasmState []string, stageStatus common.StageStatusType) ([]*models.CustodialPaymentAddressAction, error)
	UpdateQueueRunning(item *models.CustodialPaymentAddressAction) error
	UpdateQueueRetry(retryTimes int, withdraw *models.CustodialPaymentAddressAction) error
	GetListTimeout(retryTimes int) ([]*models.CustodialPaymentAddressAction, error)
}

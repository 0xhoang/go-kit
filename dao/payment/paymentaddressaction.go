package payment

import (
	"github.com/0xhoang/go-kit/models"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"time"
)

type PaymentAddressAction struct {
	db *gorm.DB
}

func NewPaymentAddressAction(db *gorm.DB) *PaymentAddressAction {
	return &PaymentAddressAction{db: db}
}

func (o *PaymentAddressAction) Update(model *models.CustodialPaymentAddressAction) error {
	if err := o.db.Save(model).Error; err != nil {
		return errors.Wrap(err, "tx.Update")
	}

	return nil
}

func (u *PaymentAddressAction) GetJobProcessing(userId uint64) ([]*models.CustodialPaymentAddressAction, error) {
	var results []*models.CustodialPaymentAddressAction

	if err := u.db.Model(&models.CustodialPaymentAddressAction{}).
		Where("aasm_state in (?) and entity_id = ? and stage_status = 0",
			[]string{
				models.AasmStateSubmitted,
				models.AasmStateProcessing,
				models.AasmStateConfirming}, userId).
		Order("id asc").Limit(30).
		Find(&results).Error; err != nil {
		return results, errors.Wrap(err, "u.Scan")
	}

	return results, nil
}

func (o *PaymentAddressAction) CreateJob(model *models.CustodialPaymentAddressAction) error {
	if err := o.db.Create(model).Error; err != nil {
		return errors.Wrap(err, "tx.Update")
	}

	return nil
}

func (u *PaymentAddressAction) UpdateRawData(id uint, rawData, txHash string) error {
	tx := u.db.Begin()
	if err := tx.Error; err != nil {
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var model *models.CustodialPaymentAddressAction
	if err := tx.Model(&model).Where("id = ?", id).Updates(map[string]interface{}{
		"tx_hash":    txHash,
		"raw_data":   rawData,
		"updated_at": time.Now().UTC(),
	}).Error; err != nil {
		tx.Rollback()
		return errors.Wrap(err, "Updates")
	}

	return tx.Commit().Error
}

func (u *PaymentAddressAction) UpdateState(id uint, aasmState string, stageStatus int, errCount int) error {
	tx := u.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}
	var model *models.CustodialPaymentAddressAction

	updateObj := map[string]interface{}{
		"aasm_state":   aasmState,
		"stage_status": stageStatus,
		"err_count":    errCount,
		"updated_at":   time.Now().UTC(),
	}

	if aasmState == models.AasmStateSucceed || aasmState == models.AasmStateErrored {
		updateObj["completed_at"] = time.Now().UTC()
	}

	if err := tx.Model(&model).Where("id = ?", id).Updates(updateObj).Error; err != nil {
		tx.Rollback()
		return errors.Wrap(err, "Updates")
	}

	return tx.Commit().Error
}

func (u *PaymentAddressAction) InsertLog(data *models.CustodialPaymentAddressLog) error {
	return errors.Wrap(u.db.Create(data).Error, "u.db.Create")
}

// start
func (o *PaymentAddressAction) ListQueue(aasmState []string, stage []string, stageStatus models.StageStatusType) ([]*models.CustodialPaymentAddressAction, error) {
	var results []*models.CustodialPaymentAddressAction

	if err := o.db.Where("aasm_state in (?) and stage_status = (?)",
		aasmState, stageStatus).
		Order("id asc").
		Limit(30).
		Find(&results).Error; err != nil {
		return nil, errors.Wrap(err, "c.db.Where")
	}

	return results, nil
}

func (o *PaymentAddressAction) UpdateQueueRunning(item *models.CustodialPaymentAddressAction) error {
	var model *models.CustodialPaymentAddressAction
	if err := o.db.Model(&model).Where("id = ?", item.ID).Updates(map[string]interface{}{
		"stage_status": models.StageStatusRunningStatus,
		"updated_at":   time.Now().UTC(),
	}).Error; err != nil {
		return errors.Wrap(err, "Updates")
	}

	return nil
}

func (o *PaymentAddressAction) UpdateQueueRetry(retryTimes int, withdraw *models.CustodialPaymentAddressAction) error {
	errCount := withdraw.ErrCount + 1

	stageStatus := models.StageStatus
	if withdraw.ErrCount > retryTimes {
		stageStatus = models.StageStatusFailedStatus
	}

	var model *models.CustodialPaymentAddressAction
	if err := o.db.Model(&model).Where("id = ?", withdraw.ID).Updates(map[string]interface{}{
		"err_count":    errCount,
		"stage_status": stageStatus,
		"updated_at":   time.Now().UTC(),
	}).Error; err != nil {
		return errors.Wrap(err, "Updates")
	}

	return nil
}

func (o *PaymentAddressAction) GetListTimeout(retryTimes int) ([]*models.CustodialPaymentAddressAction, error) {
	var results []*models.CustodialPaymentAddressAction

	ignoreAasmState := []string{
		models.AasmStateErrored,
		models.AasmStateSucceed,
	}

	//ignore: "(stage_status = ? or stage_status = ? or err_count >= ?) and aasm_state not in (?)"

	if err := o.db.Where(
		"(stage_status = ? or err_count >= ?) and aasm_state not in (?)",
		models.StageStatusRunningStatus,
		retryTimes,
		ignoreAasmState,
	).Find(&results).Limit(30).Error; err != nil {
		return nil, errors.Wrap(err, "c.db.Where")
	}

	return results, nil
}

//end

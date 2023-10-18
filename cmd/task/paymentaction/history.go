package paymentaction

import (
	"fmt"
	taskCommon "github.com/0xhoang/go-kit/cmd/task/common"
	"github.com/0xhoang/go-kit/internal/models"
	"github.com/pkg/errors"
)

type HistoryAction struct {
	model    *models.CustodialPaymentAddressLog
	property *taskCommon.AssetsProperties
}

func NewHistoryAction(model *models.CustodialPaymentAddressLog, property *taskCommon.AssetsProperties) *HistoryAction {
	return &HistoryAction{model: model, property: property}
}

func (s *HistoryAction) Task() {
	s.property.Logger.Info("Start run HistoryAction...")

	if err := s.property.Dao.InsertLog(s.model); err != nil {
		err = errors.WithStack(err)
		s.property.Logger.Error(fmt.Sprintf("Error updating log: id %v err %v", s.model.ID, err))
	}

	s.property.Logger.Info("The end HistoryAction!")
}

func (s *HistoryAction) TrackHistory(status string, message string, responseData string) {
}

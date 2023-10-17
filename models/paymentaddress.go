package models

import (
	"gorm.io/gorm"
	"time"
)

type CustodialPaymentAddressAction struct {
	gorm.Model
	EntityId    uint64 `gorm:"index:entity_id"`
	From        string
	To          string
	Amount      string
	AasmState   string          `gorm:"index:aasm_state"`
	StageStatus StageStatusType `gorm:"index:stage_status"`
	Error       string          `gorm:"type:text"`
	ErrCount    int
	CompletedAt *time.Time
}

type CustodialPaymentAddressLog struct {
	ID        uint64 `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time

	CustodialPaymentActionID uint `gorm:"index:custodial_payment_action_id"`
	Action                   string
	State                    string
	Status                   string
	Msg                      string
	Data                     string `gorm:"type:text"`
}

package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Approval struct {
	ID         uuid.UUID `gorm:"type:char(36);primaryKey" json:"id"`
	ExpenseID  uuid.UUID `gorm:"type:char(36);index;not null" json:"expense_id"`
	ApproverID uuid.UUID `gorm:"type:char(36);index;not null" json:"approver_id"`
	Status     string    `gorm:"type:varchar(30);not null" json:"status"`
	Notes      string    `gorm:"type:text" json:"notes,omitempty"`
	CreatedAt  time.Time `gorm:"column:created_at;autoCreateTime:milli" json:"created_at"`
	Expense    Expense   `gorm:"foreignKey:ExpenseID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	Approver   User      `gorm:"foreignKey:ApproverID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"-"`
}

func (a *Approval) TableName() string {
	return "approvals"
}

func (a *Approval) BeforeCreate(_ *gorm.DB) (err error) {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return
}

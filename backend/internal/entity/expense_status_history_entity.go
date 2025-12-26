package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ExpenseStatusHistory struct {
	ID             uuid.UUID  `gorm:"type:char(36);primaryKey" json:"id"`
	ExpenseID      uuid.UUID  `gorm:"type:char(36);index;not null" json:"expense_id"`
	ActorID        *uuid.UUID `gorm:"type:char(36);index" json:"actor_id,omitempty"`
	PreviousStatus string     `gorm:"type:varchar(30)" json:"previous_status,omitempty"`
	NewStatus      string     `gorm:"type:varchar(30);not null" json:"new_status"`
	Notes          string     `gorm:"type:text" json:"notes,omitempty"`
	CreatedAt      time.Time  `gorm:"column:created_at;autoCreateTime:milli" json:"created_at"`
	Expense        Expense    `gorm:"foreignKey:ExpenseID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	Actor          *User      `gorm:"foreignKey:ActorID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"-"`
}

func (e *ExpenseStatusHistory) TableName() string {
	return "expense_status_histories"
}

func (e *ExpenseStatusHistory) BeforeCreate(_ *gorm.DB) (err error) {
	if e.ID == uuid.Nil {
		e.ID = uuid.New()
	}
	return
}

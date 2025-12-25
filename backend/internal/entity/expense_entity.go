package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Expense struct {
	ID          uuid.UUID  `gorm:"type:char(36);primaryKey" json:"id"`
	UserID      uuid.UUID  `gorm:"type:char(36);index;not null" json:"user_id"`
	AmountIDR   int64      `gorm:"not null" json:"amount_idr"`
	Description string     `gorm:"type:varchar(255);not null" json:"description"`
	ReceiptURL  string     `gorm:"type:text" json:"receipt_url,omitempty"`
	Status      string     `gorm:"type:varchar(30);not null" json:"status"`
	SubmittedAt time.Time  `gorm:"column:submitted_at;autoCreateTime:milli" json:"submitted_at"`
	ProcessedAt *time.Time `gorm:"column:processed_at" json:"processed_at,omitempty"`
	CreatedAt   time.Time  `gorm:"column:created_at;autoCreateTime:milli" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"column:updated_at;autoCreateTime:milli;autoUpdateTime:milli" json:"updated_at"`
}

func (e *Expense) TableName() string {
	return "expenses"
}

func (e *Expense) BeforeCreate(_ *gorm.DB) (err error) {
	if e.ID == uuid.Nil {
		e.ID = uuid.New()
	}
	return
}

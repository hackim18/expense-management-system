package model

import (
	"time"

	"github.com/google/uuid"
)

type CreateExpenseRequest struct {
	AmountIDR   int64  `json:"amount_idr" validate:"required"`
	Description string `json:"description" validate:"required"`
	ReceiptURL  string `json:"receipt_url,omitempty" validate:"omitempty"`
}

type ExpenseResponse struct {
	ID                 uuid.UUID  `json:"id"`
	UserID             *uuid.UUID `json:"user_id,omitempty"`
	AmountIDR          int64      `json:"amount_idr"`
	AmountIDRFormatted string     `json:"amount_idr_formatted,omitempty"`
	Description        string     `json:"description"`
	ReceiptURL         string     `json:"receipt_url,omitempty"`
	Status             string     `json:"status"`
	RequiresApproval   bool       `json:"requires_approval"`
	AutoApproved       bool       `json:"auto_approved"`
	SubmittedAt        time.Time  `json:"submitted_at"`
	ProcessedAt        *time.Time `json:"processed_at,omitempty"`
}

type ExpenseDetailResponse struct {
	ExpenseResponse
	Approvals []ApprovalResponse `json:"approvals,omitempty"`
}

type ApprovalResponse struct {
	ID         uuid.UUID `json:"id"`
	ExpenseID  uuid.UUID `json:"expense_id"`
	ApproverID uuid.UUID `json:"approver_id"`
	Status     string    `json:"status"`
	Notes      string    `json:"notes,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
}

type ApproveExpenseRequest struct {
	Notes string `json:"notes,omitempty" validate:"max=500"`
}

type ExpenseStatusHistoryResponse struct {
	ID             uuid.UUID  `json:"id"`
	ExpenseID      uuid.UUID  `json:"expense_id"`
	ActorID        *uuid.UUID `json:"actor_id,omitempty"`
	PreviousStatus string     `json:"previous_status,omitempty"`
	NewStatus      string     `json:"new_status"`
	Notes          string     `json:"notes,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
}

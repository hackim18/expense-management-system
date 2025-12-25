package model

import "github.com/google/uuid"

type PaymentRequest struct {
	Amount     int64  `json:"amount"`
	ExternalID string `json:"external_id"`
}

type PaymentResponse struct {
	ID         string `json:"id"`
	ExternalID string `json:"external_id"`
	Status     string `json:"status"`
}

type PaymentJob struct {
	ExpenseID  uuid.UUID
	AmountIDR  int64
	ExternalID string
}

package usecase

import (
	"context"
	"go-expense-management-system/internal/model"
)

type PaymentProcessor interface {
	Process(ctx context.Context, request model.PaymentRequest) (*model.PaymentResponse, error)
}

type PaymentQueue interface {
	Enqueue(job model.PaymentJob) bool
}

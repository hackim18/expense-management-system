package usecase

import (
	"context"
	"go-expense-management-system/internal/model"
)

type EmailSender interface {
	Send(ctx context.Context, request model.EmailRequest) error
}

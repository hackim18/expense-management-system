package repository

import (
	"go-expense-management-system/internal/entity"

	"github.com/sirupsen/logrus"
)

type ExpenseStatusHistoryRepository struct {
	Repository[entity.ExpenseStatusHistory]
	Log *logrus.Logger
}

func NewExpenseStatusHistoryRepository(log *logrus.Logger) *ExpenseStatusHistoryRepository {
	return &ExpenseStatusHistoryRepository{
		Log: log,
	}
}

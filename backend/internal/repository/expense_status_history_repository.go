package repository

import (
	"go-expense-management-system/internal/entity"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
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

func (r *ExpenseStatusHistoryRepository) ListByExpenseID(db *gorm.DB, expenseID uuid.UUID) ([]entity.ExpenseStatusHistory, error) {
	var histories []entity.ExpenseStatusHistory
	if err := db.Where("expense_id = ?", expenseID).Order("created_at asc").Find(&histories).Error; err != nil {
		return nil, err
	}
	return histories, nil
}

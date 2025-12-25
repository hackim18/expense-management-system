package repository

import (
	"go-expense-management-system/internal/entity"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ExpenseRepository struct {
	Repository[entity.Expense]
	Log *logrus.Logger
}

type ExpenseFilter struct {
	UserID *uuid.UUID
	Status *string
}

func NewExpenseRepository(log *logrus.Logger) *ExpenseRepository {
	return &ExpenseRepository{
		Log: log,
	}
}

func (r *ExpenseRepository) List(db *gorm.DB, filter ExpenseFilter, page, size int) ([]entity.Expense, int64, error) {
	var expenses []entity.Expense
	var total int64

	query := db.Model(&entity.Expense{})
	if filter.UserID != nil {
		query = query.Where("user_id = ?", *filter.UserID)
	}
	if filter.Status != nil && *filter.Status != "" {
		query = query.Where("status = ?", *filter.Status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if size > 0 {
		offset := (page - 1) * size
		query = query.Offset(offset).Limit(size)
	}

	if err := query.Order("submitted_at desc").Find(&expenses).Error; err != nil {
		return nil, 0, err
	}

	return expenses, total, nil
}

package repository

import (
	"go-expense-management-system/internal/entity"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ApprovalRepository struct {
	Repository[entity.Approval]
	Log *logrus.Logger
}

func NewApprovalRepository(log *logrus.Logger) *ApprovalRepository {
	return &ApprovalRepository{
		Log: log,
	}
}

func (r *ApprovalRepository) ListByExpenseID(db *gorm.DB, expenseID uuid.UUID) ([]entity.Approval, error) {
	var approvals []entity.Approval
	if err := db.Where("expense_id = ?", expenseID).Order("created_at asc").Find(&approvals).Error; err != nil {
		return nil, err
	}
	return approvals, nil
}

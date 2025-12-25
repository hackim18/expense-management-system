package migrations

import (
	"go-expense-management-system/internal/entity"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&entity.User{}, &entity.Expense{}, &entity.Approval{}, &entity.ExpenseStatusHistory{})
}

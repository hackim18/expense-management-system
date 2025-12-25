package converter

import (
	"go-expense-management-system/internal/constants"
	"go-expense-management-system/internal/entity"
	"go-expense-management-system/internal/model"
	"go-expense-management-system/internal/utils"
)

func ExpenseToResponse(expense *entity.Expense, includeUserID bool) *model.ExpenseResponse {
	response := &model.ExpenseResponse{
		ID:                 expense.ID,
		AmountIDR:          expense.AmountIDR,
		AmountIDRFormatted: utils.FormatIDR(expense.AmountIDR),
		Description:        expense.Description,
		ReceiptURL:         expense.ReceiptURL,
		Status:             expense.Status,
		RequiresApproval:   expense.AmountIDR >= constants.ApprovalThreshold,
		AutoApproved:       expense.AmountIDR < constants.ApprovalThreshold && expense.Status != constants.ExpenseStatusAwaitingApproval,
		SubmittedAt:        expense.SubmittedAt,
		ProcessedAt:        expense.ProcessedAt,
	}

	if includeUserID {
		userID := expense.UserID
		response.UserID = &userID
	}

	return response
}

func ApprovalToResponse(approval *entity.Approval) model.ApprovalResponse {
	return model.ApprovalResponse{
		ID:         approval.ID,
		ExpenseID:  approval.ExpenseID,
		ApproverID: approval.ApproverID,
		Status:     approval.Status,
		Notes:      approval.Notes,
		CreatedAt:  approval.CreatedAt,
	}
}

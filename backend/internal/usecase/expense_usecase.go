package usecase

import (
	"context"
	"go-expense-management-system/internal/constants"
	"go-expense-management-system/internal/entity"
	"go-expense-management-system/internal/messages"
	"go-expense-management-system/internal/model"
	"go-expense-management-system/internal/model/converter"
	"go-expense-management-system/internal/repository"
	"go-expense-management-system/internal/utils"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ExpenseUseCase struct {
	DB                 *gorm.DB
	Log                *logrus.Logger
	ExpenseRepository  *repository.ExpenseRepository
	ApprovalRepository *repository.ApprovalRepository
	HistoryRepository  *repository.ExpenseStatusHistoryRepository
	PaymentQueue       PaymentQueue
	PaymentProcessor   PaymentProcessor
}

func NewExpenseUseCase(
	db *gorm.DB,
	logger *logrus.Logger,
	expenseRepository *repository.ExpenseRepository,
	approvalRepository *repository.ApprovalRepository,
	historyRepository *repository.ExpenseStatusHistoryRepository,
	paymentQueue PaymentQueue,
	paymentProcessor PaymentProcessor,
) *ExpenseUseCase {
	return &ExpenseUseCase{
		DB:                 db,
		Log:                logger,
		ExpenseRepository:  expenseRepository,
		ApprovalRepository: approvalRepository,
		HistoryRepository:  historyRepository,
		PaymentQueue:       paymentQueue,
		PaymentProcessor:   paymentProcessor,
	}
}

func (c *ExpenseUseCase) Create(ctx context.Context, auth *model.Auth, request *model.CreateExpenseRequest) (*model.ExpenseResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := validateExpenseAmount(request.AmountIDR); err != nil {
		return nil, err
	}

	description := strings.TrimSpace(request.Description)
	if description == "" {
		return nil, utils.Error(messages.InvalidRequestData, http.StatusBadRequest, nil)
	}

	requiresApproval := request.AmountIDR >= constants.ApprovalThreshold
	status := constants.ExpenseStatusAutoApproved
	if requiresApproval {
		status = constants.ExpenseStatusAwaitingApproval
	}

	expense := &entity.Expense{
		UserID:      auth.UserID,
		AmountIDR:   request.AmountIDR,
		Description: description,
		ReceiptURL:  request.ReceiptURL,
		Status:      status,
	}

	if err := c.ExpenseRepository.Create(tx, expense); err != nil {
		c.Log.Warnf("Failed to create expense: %+v", err)
		return nil, utils.Error(messages.InternalServerError, http.StatusInternalServerError, err)
	}
	if err := c.recordStatusChange(tx, expense, &auth.UserID, "", expense.Status, ""); err != nil {
		c.Log.Warnf("Failed to create expense history: %+v", err)
		return nil, utils.Error(messages.InternalServerError, http.StatusInternalServerError, err)
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed to commit transaction: %+v", err)
		return nil, utils.Error(messages.ErrCommitTransaction, http.StatusInternalServerError, err)
	}

	if !requiresApproval {
		c.enqueuePayment(expense)
	}

	return converter.ExpenseToResponse(expense, false), nil
}

func (c *ExpenseUseCase) List(ctx context.Context, auth *model.Auth, status string, page, size int) ([]model.ExpenseResponse, model.PageMetadata, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	filter := repository.ExpenseFilter{}
	if !isManager(auth) {
		filter.UserID = &auth.UserID
	}

	status = normalizeStatusFilter(status)
	if status != "" {
		filter.Status = &status
	}

	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 10
	}

	expenses, total, err := c.ExpenseRepository.List(tx, filter, page, size)
	if err != nil {
		c.Log.Warnf("Failed to list expenses: %+v", err)
		return nil, model.PageMetadata{}, utils.Error(messages.InternalServerError, http.StatusInternalServerError, err)
	}

	responses := make([]model.ExpenseResponse, 0, len(expenses))
	includeUserID := isManager(auth)
	for i := range expenses {
		responses = append(responses, *converter.ExpenseToResponse(&expenses[i], includeUserID))
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed to commit transaction: %+v", err)
		return nil, model.PageMetadata{}, utils.Error(messages.ErrCommitTransaction, http.StatusInternalServerError, err)
	}

	paging := utils.NewPageMetadata(page, size, total)
	return responses, paging, nil
}

func (c *ExpenseUseCase) Get(ctx context.Context, auth *model.Auth, expenseID uuid.UUID) (*model.ExpenseDetailResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	expense := new(entity.Expense)
	if err := c.ExpenseRepository.FindById(tx, expense, expenseID); err != nil {
		return nil, utils.Error(messages.ErrExpenseNotFound, http.StatusNotFound, err)
	}

	if !isManager(auth) && expense.UserID != auth.UserID {
		return nil, utils.Error(messages.Forbidden, http.StatusForbidden, nil)
	}

	approvals, err := c.ApprovalRepository.ListByExpenseID(tx, expense.ID)
	if err != nil {
		c.Log.Warnf("Failed to list approvals: %+v", err)
		return nil, utils.Error(messages.InternalServerError, http.StatusInternalServerError, err)
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed to commit transaction: %+v", err)
		return nil, utils.Error(messages.ErrCommitTransaction, http.StatusInternalServerError, err)
	}

	includeUserID := isManager(auth)
	response := model.ExpenseDetailResponse{
		ExpenseResponse: *converter.ExpenseToResponse(expense, includeUserID),
	}
	if len(approvals) > 0 {
		response.Approvals = make([]model.ApprovalResponse, 0, len(approvals))
		for i := range approvals {
			response.Approvals = append(response.Approvals, converter.ApprovalToResponse(&approvals[i]))
		}
	}

	return &response, nil
}

func (c *ExpenseUseCase) History(ctx context.Context, auth *model.Auth, expenseID uuid.UUID) ([]model.ExpenseStatusHistoryResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	expense := new(entity.Expense)
	if err := c.ExpenseRepository.FindById(tx, expense, expenseID); err != nil {
		return nil, utils.Error(messages.ErrExpenseNotFound, http.StatusNotFound, err)
	}

	if !isManager(auth) && expense.UserID != auth.UserID {
		return nil, utils.Error(messages.Forbidden, http.StatusForbidden, nil)
	}

	histories, err := c.HistoryRepository.ListByExpenseID(tx, expense.ID)
	if err != nil {
		c.Log.Warnf("Failed to list expense histories: %+v", err)
		return nil, utils.Error(messages.InternalServerError, http.StatusInternalServerError, err)
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed to commit transaction: %+v", err)
		return nil, utils.Error(messages.ErrCommitTransaction, http.StatusInternalServerError, err)
	}

	responses := make([]model.ExpenseStatusHistoryResponse, 0, len(histories))
	for i := range histories {
		responses = append(responses, converter.ExpenseStatusHistoryToResponse(&histories[i]))
	}

	return responses, nil
}

func (c *ExpenseUseCase) Approve(ctx context.Context, auth *model.Auth, expenseID uuid.UUID, request *model.ApproveExpenseRequest) (*model.ExpenseResponse, error) {
	if !isManager(auth) {
		return nil, utils.Error(messages.Forbidden, http.StatusForbidden, nil)
	}

	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	expense := new(entity.Expense)
	if err := c.ExpenseRepository.FindById(tx, expense, expenseID); err != nil {
		return nil, utils.Error(messages.ErrExpenseNotFound, http.StatusNotFound, err)
	}

	if expense.Status != constants.ExpenseStatusAwaitingApproval {
		return nil, utils.Error(messages.ErrExpenseNotPending, http.StatusConflict, nil)
	}

	approval := &entity.Approval{
		ExpenseID:  expense.ID,
		ApproverID: auth.UserID,
		Status:     constants.ApprovalStatusApproved,
		Notes:      strings.TrimSpace(request.Notes),
	}

	if err := c.ApprovalRepository.Create(tx, approval); err != nil {
		c.Log.Warnf("Failed to create approval: %+v", err)
		return nil, utils.Error(messages.InternalServerError, http.StatusInternalServerError, err)
	}

	previousStatus := expense.Status
	expense.Status = constants.ExpenseStatusApproved
	if err := c.ExpenseRepository.Update(tx, expense); err != nil {
		c.Log.Warnf("Failed to update expense: %+v", err)
		return nil, utils.Error(messages.InternalServerError, http.StatusInternalServerError, err)
	}
	if err := c.recordStatusChange(tx, expense, &auth.UserID, previousStatus, expense.Status, approval.Notes); err != nil {
		c.Log.Warnf("Failed to create expense history: %+v", err)
		return nil, utils.Error(messages.InternalServerError, http.StatusInternalServerError, err)
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed to commit transaction: %+v", err)
		return nil, utils.Error(messages.ErrCommitTransaction, http.StatusInternalServerError, err)
	}

	c.enqueuePayment(expense)
	return converter.ExpenseToResponse(expense, true), nil
}

func (c *ExpenseUseCase) Reject(ctx context.Context, auth *model.Auth, expenseID uuid.UUID, request *model.ApproveExpenseRequest) (*model.ExpenseResponse, error) {
	if !isManager(auth) {
		return nil, utils.Error(messages.Forbidden, http.StatusForbidden, nil)
	}

	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	expense := new(entity.Expense)
	if err := c.ExpenseRepository.FindById(tx, expense, expenseID); err != nil {
		return nil, utils.Error(messages.ErrExpenseNotFound, http.StatusNotFound, err)
	}

	if expense.Status != constants.ExpenseStatusAwaitingApproval {
		return nil, utils.Error(messages.ErrExpenseNotPending, http.StatusConflict, nil)
	}

	approval := &entity.Approval{
		ExpenseID:  expense.ID,
		ApproverID: auth.UserID,
		Status:     constants.ApprovalStatusRejected,
		Notes:      strings.TrimSpace(request.Notes),
	}

	if err := c.ApprovalRepository.Create(tx, approval); err != nil {
		c.Log.Warnf("Failed to create approval: %+v", err)
		return nil, utils.Error(messages.InternalServerError, http.StatusInternalServerError, err)
	}

	previousStatus := expense.Status
	expense.Status = constants.ExpenseStatusRejected
	if err := c.ExpenseRepository.Update(tx, expense); err != nil {
		c.Log.Warnf("Failed to update expense: %+v", err)
		return nil, utils.Error(messages.InternalServerError, http.StatusInternalServerError, err)
	}
	if err := c.recordStatusChange(tx, expense, &auth.UserID, previousStatus, expense.Status, approval.Notes); err != nil {
		c.Log.Warnf("Failed to create expense history: %+v", err)
		return nil, utils.Error(messages.InternalServerError, http.StatusInternalServerError, err)
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed to commit transaction: %+v", err)
		return nil, utils.Error(messages.ErrCommitTransaction, http.StatusInternalServerError, err)
	}

	return converter.ExpenseToResponse(expense, true), nil
}

func (c *ExpenseUseCase) ProcessPayment(ctx context.Context, job model.PaymentJob) error {
	if c.PaymentProcessor == nil {
		return utils.Error(messages.ErrPaymentFailed, http.StatusBadGateway, nil)
	}

	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	expense := new(entity.Expense)
	if err := tx.Where("id = ?", job.ExpenseID).Take(expense).Error; err != nil {
		return utils.Error(messages.ErrExpenseNotFound, http.StatusNotFound, err)
	}

	if expense.Status == constants.ExpenseStatusCompleted || expense.Status == constants.ExpenseStatusRejected {
		return nil
	}

	if expense.Status != constants.ExpenseStatusApproved && expense.Status != constants.ExpenseStatusAutoApproved {
		return nil
	}

	if expense.ProcessedAt != nil {
		return nil
	}

	_, err := c.PaymentProcessor.Process(ctx, model.PaymentRequest{
		Amount:     job.AmountIDR,
		ExternalID: job.ExternalID,
	})
	if err != nil {
		c.Log.Warnf("Payment processing failed for expense %s: %+v", job.ExpenseID, err)
		return utils.Error(messages.ErrPaymentFailed, http.StatusBadGateway, err)
	}

	now := time.Now()
	previousStatus := expense.Status
	expense.Status = constants.ExpenseStatusCompleted
	expense.ProcessedAt = &now
	if err := tx.Save(expense).Error; err != nil {
		c.Log.Warnf("Failed to update expense payment status: %+v", err)
		return utils.Error(messages.InternalServerError, http.StatusInternalServerError, err)
	}
	if err := c.recordStatusChange(tx, expense, nil, previousStatus, expense.Status, ""); err != nil {
		c.Log.Warnf("Failed to create expense history: %+v", err)
		return utils.Error(messages.InternalServerError, http.StatusInternalServerError, err)
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed to commit transaction: %+v", err)
		return utils.Error(messages.ErrCommitTransaction, http.StatusInternalServerError, err)
	}

	return nil
}

func validateExpenseAmount(amount int64) error {
	if amount <= 0 {
		return utils.Error(messages.ErrInvalidExpenseAmount, http.StatusBadRequest, nil)
	}
	if amount < constants.MinExpenseAmount || amount > constants.MaxExpenseAmount {
		return utils.Error(messages.ErrInvalidExpenseAmount, http.StatusBadRequest, nil)
	}
	return nil
}

func normalizeStatusFilter(status string) string {
	status = strings.TrimSpace(strings.ToLower(status))
	switch status {
	case "pending":
		return constants.ExpenseStatusAwaitingApproval
	case "auto-approved":
		return constants.ExpenseStatusAutoApproved
	}
	return status
}

func (c *ExpenseUseCase) enqueuePayment(expense *entity.Expense) {
	if c.PaymentQueue == nil {
		return
	}

	job := model.PaymentJob{
		ExpenseID:  expense.ID,
		AmountIDR:  expense.AmountIDR,
		ExternalID: expense.ID.String(),
	}
	c.PaymentQueue.Enqueue(job)
}

func isManager(auth *model.Auth) bool {
	if auth == nil {
		return false
	}
	return auth.Role == constants.RoleManager
}

func (c *ExpenseUseCase) recordStatusChange(
	tx *gorm.DB,
	expense *entity.Expense,
	actorID *uuid.UUID,
	previousStatus string,
	newStatus string,
	notes string,
) error {
	if c.HistoryRepository == nil {
		return nil
	}

	history := &entity.ExpenseStatusHistory{
		ExpenseID:      expense.ID,
		ActorID:        actorID,
		PreviousStatus: previousStatus,
		NewStatus:      newStatus,
		Notes:          strings.TrimSpace(notes),
	}

	return c.HistoryRepository.Create(tx, history)
}

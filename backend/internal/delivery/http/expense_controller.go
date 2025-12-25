package http

import (
	"errors"
	"go-expense-management-system/internal/delivery/http/middleware"
	"go-expense-management-system/internal/messages"
	"go-expense-management-system/internal/model"
	"go-expense-management-system/internal/usecase"
	"go-expense-management-system/internal/utils"
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type ExpenseController struct {
	Log      *logrus.Logger
	UseCase  *usecase.ExpenseUseCase
	Validate *validator.Validate
}

func NewExpenseController(useCase *usecase.ExpenseUseCase, logger *logrus.Logger, validate *validator.Validate) *ExpenseController {
	return &ExpenseController{
		Log:      logger,
		UseCase:  useCase,
		Validate: validate,
	}
}

func (c *ExpenseController) Create(ctx *gin.Context) {
	auth, ok := getAuthOrAbort(ctx)
	if !ok {
		return
	}

	request := new(model.CreateExpenseRequest)
	if err := ctx.ShouldBindJSON(request); err != nil {
		c.Log.Warnf("Failed to parse request body: %+v", err)
		utils.HandleHTTPError(ctx, utils.Error(messages.FailedDataFromBody, http.StatusBadRequest, err))
		return
	}

	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("Validation failed: %+v", err)
		message := utils.TranslateValidationError(c.Validate, err)
		utils.HandleHTTPError(ctx, utils.Error(message, http.StatusBadRequest, err))
		return
	}

	response, err := c.UseCase.Create(ctx.Request.Context(), auth, request)
	if err != nil {
		c.Log.Warnf("Failed to create expense: %+v", err)
		utils.HandleHTTPError(ctx, err)
		return
	}

	res := utils.SuccessResponse(messages.ExpenseCreated, response)
	ctx.JSON(http.StatusCreated, res)
}

func (c *ExpenseController) List(ctx *gin.Context) {
	auth, ok := getAuthOrAbort(ctx)
	if !ok {
		return
	}

	status := ctx.Query("status")
	page := parseIntQuery(ctx.Query("page"), 1)
	size := parseIntQuery(ctx.Query("size"), 10)

	responses, paging, err := c.UseCase.List(ctx.Request.Context(), auth, status, page, size)
	if err != nil {
		c.Log.Warnf("Failed to list expenses: %+v", err)
		utils.HandleHTTPError(ctx, err)
		return
	}

	res := utils.SuccessWithPaginationResponse(messages.ExpenseListed, responses, paging)
	ctx.JSON(http.StatusOK, res)
}

func (c *ExpenseController) Get(ctx *gin.Context) {
	auth, ok := getAuthOrAbort(ctx)
	if !ok {
		return
	}

	expenseID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.HandleHTTPError(ctx, utils.Error(messages.ErrInvalidIDFormat, http.StatusBadRequest, err))
		return
	}

	response, err := c.UseCase.Get(ctx.Request.Context(), auth, expenseID)
	if err != nil {
		c.Log.Warnf("Failed to fetch expense: %+v", err)
		utils.HandleHTTPError(ctx, err)
		return
	}

	res := utils.SuccessResponse(messages.ExpenseFetched, response)
	ctx.JSON(http.StatusOK, res)
}

func (c *ExpenseController) Approve(ctx *gin.Context) {
	auth, ok := getAuthOrAbort(ctx)
	if !ok {
		return
	}

	expenseID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.HandleHTTPError(ctx, utils.Error(messages.ErrInvalidIDFormat, http.StatusBadRequest, err))
		return
	}

	request := new(model.ApproveExpenseRequest)
	if err := ctx.ShouldBindJSON(request); err != nil && !errors.Is(err, io.EOF) {
		c.Log.Warnf("Failed to parse request body: %+v", err)
		utils.HandleHTTPError(ctx, utils.Error(messages.FailedDataFromBody, http.StatusBadRequest, err))
		return
	}

	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("Validation failed: %+v", err)
		message := utils.TranslateValidationError(c.Validate, err)
		utils.HandleHTTPError(ctx, utils.Error(message, http.StatusBadRequest, err))
		return
	}

	response, err := c.UseCase.Approve(ctx.Request.Context(), auth, expenseID, request)
	if err != nil {
		c.Log.Warnf("Failed to approve expense: %+v", err)
		utils.HandleHTTPError(ctx, err)
		return
	}

	res := utils.SuccessResponse(messages.ExpenseApproved, response)
	ctx.JSON(http.StatusOK, res)
}

func (c *ExpenseController) Reject(ctx *gin.Context) {
	auth, ok := getAuthOrAbort(ctx)
	if !ok {
		return
	}

	expenseID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.HandleHTTPError(ctx, utils.Error(messages.ErrInvalidIDFormat, http.StatusBadRequest, err))
		return
	}

	request := new(model.ApproveExpenseRequest)
	if err := ctx.ShouldBindJSON(request); err != nil && !errors.Is(err, io.EOF) {
		c.Log.Warnf("Failed to parse request body: %+v", err)
		utils.HandleHTTPError(ctx, utils.Error(messages.FailedDataFromBody, http.StatusBadRequest, err))
		return
	}

	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("Validation failed: %+v", err)
		message := utils.TranslateValidationError(c.Validate, err)
		utils.HandleHTTPError(ctx, utils.Error(message, http.StatusBadRequest, err))
		return
	}

	response, err := c.UseCase.Reject(ctx.Request.Context(), auth, expenseID, request)
	if err != nil {
		c.Log.Warnf("Failed to reject expense: %+v", err)
		utils.HandleHTTPError(ctx, err)
		return
	}

	res := utils.SuccessResponse(messages.ExpenseRejected, response)
	ctx.JSON(http.StatusOK, res)
}

func getAuthOrAbort(ctx *gin.Context) (*model.Auth, bool) {
	auth, ok := middleware.GetUser(ctx)
	if !ok {
		res := utils.FailedResponse(messages.Unauthorized)
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, res)
		return nil, false
	}
	return auth, true
}

func parseIntQuery(value string, fallback int) int {
	if value == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(value)
	if err != nil || parsed <= 0 {
		return fallback
	}
	return parsed
}

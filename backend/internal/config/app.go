package config

import (
	"go-expense-management-system/internal/background"
	"go-expense-management-system/internal/delivery/http"
	"go-expense-management-system/internal/delivery/http/middleware"
	"go-expense-management-system/internal/delivery/http/route"
	"go-expense-management-system/internal/integration/email"
	"go-expense-management-system/internal/integration/payment"
	"go-expense-management-system/internal/repository"
	"go-expense-management-system/internal/usecase"
	"go-expense-management-system/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type BootstrapConfig struct {
	Router   *gin.Engine
	DB       *gorm.DB
	JWT      *utils.JWTHelper
	Log      *logrus.Logger
	Validate *validator.Validate
	Config   *viper.Viper
}

func Bootstrap(config *BootstrapConfig) {
	// Setup repositories
	userRepository := repository.NewUserRepository(config.Log)
	expenseRepository := repository.NewExpenseRepository(config.Log)
	approvalRepository := repository.NewApprovalRepository(config.Log)
	historyRepository := repository.NewExpenseStatusHistoryRepository(config.Log)

	// Setup integrations
	paymentCfg := buildPaymentConfig(config.Config)
	paymentClient := payment.NewClient(paymentCfg.BaseURL, paymentCfg.Timeout, config.Log)

	emailClient := email.NewClient(buildSMTPConfig(config.Config), config.Log)

	// Setup use cases
	userUseCase := usecase.NewUserUseCase(config.DB, config.Log, config.JWT, userRepository)
	expenseUseCase := usecase.NewExpenseUseCase(
		config.DB,
		config.Log,
		expenseRepository,
		approvalRepository,
		historyRepository,
		userRepository,
		emailClient,
		nil,
		paymentClient,
	)

	// Setup controllers
	userController := http.NewUserController(userUseCase, config.Log, config.Validate)
	expenseController := http.NewExpenseController(expenseUseCase, config.Log, config.Validate)

	// Setup middleware
	authMiddleware := middleware.NewAuth(userUseCase)

	// Setup background workers
	paymentWorker := background.NewPaymentWorker(
		paymentCfg.QueueBuffer,
		paymentCfg.RetryCount,
		paymentCfg.RetryDelay,
		paymentCfg.Timeout,
		config.Log,
		expenseUseCase.ProcessPayment,
	)
	paymentWorker.Start()
	expenseUseCase.PaymentQueue = paymentWorker

	// Setup routes
	routeConfig := route.RouteConfig{
		Router:            config.Router,
		UserController:    userController,
		ExpenseController: expenseController,
		AuthMiddleware:    authMiddleware,
	}
	routeConfig.Setup()
}

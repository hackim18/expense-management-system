package config

import (
	"go-expense-management-system/internal/background"
	"go-expense-management-system/internal/delivery/http"
	"go-expense-management-system/internal/delivery/http/middleware"
	"go-expense-management-system/internal/delivery/http/route"
	"go-expense-management-system/internal/integration/payment"
	"go-expense-management-system/internal/repository"
	"go-expense-management-system/internal/usecase"
	"go-expense-management-system/internal/utils"
	"time"

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
	// setup repositories
	userRepository := repository.NewUserRepository(config.Log)
	expenseRepository := repository.NewExpenseRepository(config.Log)
	approvalRepository := repository.NewApprovalRepository(config.Log)

	paymentBaseURL := config.Config.GetString("PAYMENT_BASE_URL")
	paymentTimeout := time.Duration(config.Config.GetInt("PAYMENT_TIMEOUT_SECONDS")) * time.Second
	paymentRetryCount := config.Config.GetInt("PAYMENT_RETRY_COUNT")
	paymentRetryDelay := time.Duration(config.Config.GetInt("PAYMENT_RETRY_DELAY_SECONDS")) * time.Second
	paymentQueueBuffer := config.Config.GetInt("PAYMENT_QUEUE_BUFFER")

	paymentClient := payment.NewClient(paymentBaseURL, paymentTimeout, config.Log)

	// setup use cases
	userUseCase := usecase.NewUserUseCase(config.DB, config.Log, config.JWT, userRepository)
	expenseUseCase := usecase.NewExpenseUseCase(
		config.DB,
		config.Log,
		expenseRepository,
		approvalRepository,
		nil,
		paymentClient,
	)

	// setup controller
	userController := http.NewUserController(userUseCase, config.Log, config.Validate)
	expenseController := http.NewExpenseController(expenseUseCase, config.Log, config.Validate)

	// setup middleware
	authMiddleware := middleware.NewAuth(userUseCase)

	paymentWorker := background.NewPaymentWorker(
		paymentQueueBuffer,
		paymentRetryCount,
		paymentRetryDelay,
		paymentTimeout,
		config.Log,
		expenseUseCase.ProcessPayment,
	)
	paymentWorker.Start()
	expenseUseCase.PaymentQueue = paymentWorker

	routeConfig := route.RouteConfig{
		Router:            config.Router,
		UserController:    userController,
		ExpenseController: expenseController,
		AuthMiddleware:    authMiddleware,
	}
	routeConfig.Setup()
}

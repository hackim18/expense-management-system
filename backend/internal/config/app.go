package config

import (
	"go-expense-management-system/internal/delivery/http"
	"go-expense-management-system/internal/delivery/http/middleware"
	"go-expense-management-system/internal/delivery/http/route"
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
	// setup repositories
	userRepository := repository.NewUserRepository(config.Log)

	// setup use cases
	userUseCase := usecase.NewUserUseCase(config.DB, config.Log, config.JWT, userRepository)

	// setup controller
	userController := http.NewUserController(userUseCase, config.Log, config.Validate)

	// setup middleware
	authMiddleware := middleware.NewAuth(userUseCase)

	routeConfig := route.RouteConfig{
		Router:         config.Router,
		UserController: userController,
		AuthMiddleware: authMiddleware,
	}
	routeConfig.Setup()
}

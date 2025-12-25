package config

import (
	"go-expense-management-system/internal/delivery/http/middleware"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func NewGin(v *viper.Viper) *gin.Engine {
	engine := gin.New()
	engine.Use(gin.Logger(), gin.Recovery(), NewCORS(v), middleware.NewRateLimiter(v))
	engine.SetTrustedProxies(nil)
	return engine
}

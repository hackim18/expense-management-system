package route

import (
	"go-expense-management-system/internal/messages"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func (c *RouteConfig) RegisterCommonRoutes(app *gin.Engine) {
	startTime := time.Now()
	welcomeHandler := func(ctx *gin.Context) {
		res := gin.H{"message": messages.WelcomeMessage}
		ctx.JSON(http.StatusOK, res)
	}

	healthHandler := func(ctx *gin.Context) {
		res := gin.H{"status": "ok"}
		ctx.JSON(http.StatusOK, res)
	}

	metricsHandler := func(ctx *gin.Context) {
		uptime := time.Since(startTime).String()
		res := gin.H{"status": "ok", "uptime": uptime}
		ctx.JSON(http.StatusOK, res)
	}

	app.GET("/", welcomeHandler)
	app.GET("/api", welcomeHandler)
	app.GET("/health", healthHandler)
	app.GET("/api/health", healthHandler)
	app.GET("/api/metrics", metricsHandler)
	app.NoRoute(func(ctx *gin.Context) {
		res := gin.H{"message": messages.NotFound}
		ctx.AbortWithStatusJSON(http.StatusNotFound, res)
	})
}

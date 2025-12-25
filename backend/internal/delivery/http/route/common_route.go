package route

import (
	"go-expense-management-system/internal/messages"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func (c *RouteConfig) RegisterApiRoutes(api *gin.RouterGroup) {
	startTime := time.Now()
	openAPIPath := "api/openapi.yaml"

	api.GET("", c.welcomeHandler())
	api.GET("/health", c.healthHandler())
	api.GET("/openapi.yaml", func(ctx *gin.Context) {
		ctx.File(openAPIPath)
	})
	api.GET("/metrics", c.metricsHandler(startTime))
}

func (c *RouteConfig) RegisterPublicRoutes() {
	app := c.Router
	app.GET("/", c.welcomeHandler())
	app.GET("/health", c.healthHandler())
	app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("/api/openapi.yaml")))
	app.NoRoute(func(ctx *gin.Context) {
		res := gin.H{"message": messages.NotFound}
		ctx.AbortWithStatusJSON(http.StatusNotFound, res)
	})
}

func (c *RouteConfig) welcomeHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		res := gin.H{"message": messages.WelcomeMessage}
		ctx.JSON(http.StatusOK, res)
	}
}

func (c *RouteConfig) healthHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		res := gin.H{"status": "ok"}
		ctx.JSON(http.StatusOK, res)
	}
}

func (c *RouteConfig) metricsHandler(startTime time.Time) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		uptime := time.Since(startTime).String()
		res := gin.H{"status": "ok", "uptime": uptime}
		ctx.JSON(http.StatusOK, res)
	}
}

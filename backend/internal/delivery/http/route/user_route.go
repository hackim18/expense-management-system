package route

import "github.com/gin-gonic/gin"

func (c *RouteConfig) RegisterAuthRoutes(rg *gin.RouterGroup) {
	auth := rg.Group("/auth")

	auth.POST("/register", c.UserController.Register)
	auth.POST("/login", c.UserController.Login)
}

package route

import (
	"go-expense-management-system/internal/delivery/http"

	"github.com/gin-gonic/gin"
)

type RouteConfig struct {
	Router            *gin.Engine
	UserController    *http.UserController
	ExpenseController *http.ExpenseController
	AuthMiddleware    gin.HandlerFunc
}

func (c *RouteConfig) Setup() {
	api := c.Router.Group("/api")

	c.RegisterAuthRoutes(api)
	c.RegisterExpenseRoutes(api)
	c.RegisterCommonRoutes(c.Router)
}

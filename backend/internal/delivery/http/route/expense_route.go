package route

import "github.com/gin-gonic/gin"

func (c *RouteConfig) RegisterExpenseRoutes(rg *gin.RouterGroup) {
	expense := rg.Group("/expenses")
	expense.Use(c.AuthMiddleware)

	expense.POST("", c.ExpenseController.Create)
	expense.GET("", c.ExpenseController.List)
	expense.GET("/:id", c.ExpenseController.Get)
	expense.PUT("/:id/approve", c.ExpenseController.Approve)
	expense.PUT("/:id/reject", c.ExpenseController.Reject)
}

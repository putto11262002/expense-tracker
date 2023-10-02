package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/putto11262002/expense-tracker/api/handlers"
	"github.com/putto11262002/expense-tracker/api/middlewares"
	"github.com/putto11262002/expense-tracker/api/services"

)

func NewExpenseRoutes(router *gin.RouterGroup, expenseService services.IExpenseService, userService services.IUserService) {
	
	expenseHandler := handlers.NewExpenseHandler(expenseService, userService)

	expenseRG := router.Group("/expense")
	expenseRG.Use(middlewares.JWTAuthMiddleware())
	expenseRG.POST("", expenseHandler.HandleCreateExpense)
	expenseRG.GET(":id", expenseHandler.HandleGetExpenseByID)
	expenseRG.GET("", expenseHandler.HandleGetExpense)
	expenseRG.POST("dept/settle", expenseHandler.HandleSettleDepth)

}

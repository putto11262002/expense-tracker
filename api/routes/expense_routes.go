package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/putto11262002/expense-tracker/api/handlers"
	"github.com/putto11262002/expense-tracker/api/middlewares"
	"github.com/putto11262002/expense-tracker/api/repositories"
	"github.com/putto11262002/expense-tracker/api/services"
	"gorm.io/gorm"
)

func NewExpenseRoutes(r *gin.Engine, db *gorm.DB) {
	expenseRepository := repositories.NewExpenseRepository(db)
	expenseService := services.NewExpenseService(expenseRepository)
	userRepository := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepository)
	expenseHandler := handlers.NewExpenseHandler(expenseService, userService)

	expenseRG := r.Group("/expense")
	expenseRG.Use(middlewares.JWTAuthMiddleware())
	expenseRG.POST("", expenseHandler.HandleCreateExpense)
	expenseRG.GET(":id", expenseHandler.HandleGetExpenseByID)
	expenseRG.GET("", expenseHandler.HandleGetExpense)
	expenseRG.POST("dept/settle", expenseHandler.HandleSettleDepth)

}
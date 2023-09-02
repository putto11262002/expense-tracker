package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/putto11262002/expense-tracker/api/handlers"
	"github.com/putto11262002/expense-tracker/api/repositories"
	"github.com/putto11262002/expense-tracker/api/services"
	"gorm.io/gorm"
)

func NewAuthRoutes(db *gorm.DB, r *gin.Engine) {
	userRepository := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepository)
	userHandler := handlers.NewUserHandler(userService)

	userAuthRg := r.Group("/auth")

	userAuthRg.POST("/register", userHandler.HandleRegister)
	userAuthRg.POST("/login", userHandler.HandleLogin)

}

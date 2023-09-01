package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/putto11262002/expense-tracker/api/internal/handlers"
	"github.com/putto11262002/expense-tracker/api/internal/repositories"
	"github.com/putto11262002/expense-tracker/api/internal/services"
	"gorm.io/gorm"
)




func NewUserRoutes(db *gorm.DB, r *gin.Engine){
	userRepository := repositories.NewUserRepository(db)
	userService :=  services.NewUserService(userRepository)
	userHandler := handlers.NewUserHandler(userService)


	userAuthRg := r.Group("/auth")

	userAuthRg.POST("/register", userHandler.HandleRegister)


	
}





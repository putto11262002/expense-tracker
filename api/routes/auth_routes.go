package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/putto11262002/expense-tracker/api/handlers"
	"github.com/putto11262002/expense-tracker/api/services"

)

func NewAuthRoutes(router *gin.RouterGroup, userService services.IUserService) {

	userHandler := handlers.NewUserHandler(userService)

	userAuthRg := router.Group("/auth")

	userAuthRg.POST("/register", userHandler.HandleRegister)
	userAuthRg.POST("/login", userHandler.HandleLogin)

}

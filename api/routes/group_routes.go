package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/putto11262002/expense-tracker/api/handlers"
	"github.com/putto11262002/expense-tracker/api/middlewares"
	"github.com/putto11262002/expense-tracker/api/repositories"
	"github.com/putto11262002/expense-tracker/api/services"
	"gorm.io/gorm"
)

func NewGroupRoutes(db *gorm.DB, r *gin.Engine) {

	userRepository := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepository)
	groupRepository := repositories.NewGroupRepository(db)
	groupService := services.NewGroupService(groupRepository)

	groupHandlers := handlers.NewGroupHandler(groupService, userService)

	rg := r.Group("/group")

	rg.Use(middlewares.JWTAuthMiddleware())

	rg.POST("/", groupHandlers.HandleCreateGroup)
}

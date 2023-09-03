package routes

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/putto11262002/expense-tracker/api/handlers"
	"github.com/putto11262002/expense-tracker/api/middlewares"
	"github.com/putto11262002/expense-tracker/api/repositories"
	"github.com/putto11262002/expense-tracker/api/services"
	"github.com/putto11262002/expense-tracker/api/utils"
	"gorm.io/gorm"
)

func NewGroupRoutes(db *gorm.DB, r *gin.Engine) {

	userRepository := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepository)
	groupRepository := repositories.NewGroupRepository(db)
	groupService := services.NewGroupService(groupRepository)

	groupHandlers := handlers.NewGroupHandler(groupService, userService)

	rg := r.Group("/")

	rg.Use(middlewares.JWTAuthMiddleware())

	rg.POST("/group", groupHandlers.HandleCreateGroup)
	rg.GET("/user/me/group", func(ctx *gin.Context) {
		claims := utils.GetClaimsFromCtx(ctx)
		if claims == nil {
			utils.AbortWithError(ctx, errors.New("cannot retrieve claims from context"))
			return
		}
		ctx.AddParam("id", claims.Subject)
		ctx.Next()

	}, groupHandlers.HandleGetGroupsByUserID)
	rg.GET("/group/:id", groupHandlers.HandleGetGroupByID)
}

package routes

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/putto11262002/expense-tracker/api/handlers"
	"github.com/putto11262002/expense-tracker/api/middlewares"
	"github.com/putto11262002/expense-tracker/api/services"
	"github.com/putto11262002/expense-tracker/api/utils"

)

func NewGroupRoutes(router *gin.RouterGroup, groupService services.IGroupService, userService services.IUserService) {

	

	groupHandlers := handlers.NewGroupHandler(groupService, userService)

	rg := router.Group("/")

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
	rg.POST("/group/user/add", groupHandlers.HandleAddGroupMember)
}

package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/putto11262002/expense-tracker/api/handlers"
	"github.com/putto11262002/expense-tracker/api/middlewares"
	"github.com/putto11262002/expense-tracker/api/services"
	"github.com/putto11262002/expense-tracker/api/utils"

)

func NewUserRoutes(router *gin.RouterGroup, userService services.IUserService) {

	userHandler := handlers.NewUserHandler(userService)

	userRg := router.Group("/user")

	userRg.Use(middlewares.JWTAuthMiddleware())

	userRg.GET("/me",
		// retrieve user id from claims and set it as path param
		func(ctx *gin.Context) {
			val, exist := ctx.Get("claims")
			if !exist {
				utils.AbortWithError(ctx, &utils.AuthorizationError{})
				return
			}
			if claims, ok := val.(*jwt.StandardClaims); !ok {
				utils.AbortWithError(ctx, &utils.AuthorizationError{})
				return
			} else {
				ctx.AddParam("id", claims.Subject)
				ctx.Next()
			}

		}, userHandler.HandleGetUserByID)

	userRg.GET("/:id", userHandler.HandleGetUserByID)
	userRg.GET("/email/:email", userHandler.HandleUserByEmail)
	userRg.GET("", userHandler.HandleGetUsers)

}

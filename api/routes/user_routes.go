package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/putto11262002/expense-tracker/api/handlers"
	"github.com/putto11262002/expense-tracker/api/middlewares"
	"github.com/putto11262002/expense-tracker/api/repositories"
	"github.com/putto11262002/expense-tracker/api/services"
	"github.com/putto11262002/expense-tracker/api/utils"
	"gorm.io/gorm"
)

func NewUserRoutes(db *gorm.DB, r *gin.Engine) {
	userRepository := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepository)
	userHandler := handlers.NewUserHandler(userService)

	userRg := r.Group("/user")

	userRg.Use(middlewares.JWTAuthMiddleware())

	userRg.GET("/me", func(ctx *gin.Context) {
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

}

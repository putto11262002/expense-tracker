package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/putto11262002/expense-tracker/api/utils"
)

// JWTAuthMiddleware retrieves token from the request cookie
// if the token exist use the subject the claim to retrieve the associated user from the database
// else abort request with AuthenticationError
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token, err := ctx.Cookie("token")
		if err != nil {
			utils.AbortWithError(ctx, &utils.AuthorizationError{
				Message: "Invalid token",
			})
			return
		}
		claims, err := utils.ValidateToken(token, utils.GetJWTSecret())
		if err != nil {
			utils.AbortWithError(ctx, err)
			return
		}
		ctx.Set("claims", claims)
		ctx.Next()
	}
}

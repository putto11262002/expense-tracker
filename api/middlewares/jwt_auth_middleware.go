package middlewares

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/putto11262002/expense-tracker/api/utils"
)

// JWTAuthMiddleware retrieves token from the request cookie or authorization header
// if the token exist use the subject the claim to retrieve the associated user from the database
// else abort request with AuthenticationError
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token, err := ctx.Cookie("token")


		// If the token is not present in the cookie check if token is provided as a Bearer otken in the authorization header
		if token == "" || err != nil {
			authHeader := ctx.Request.Header.Get("Authorization") 
			var found bool
			token, found = strings.CutPrefix(authHeader, "Bearer ")
			if !found || token == "" {
				utils.AbortWithError(ctx, &utils.AuthorizationError{
					Message: "Invalid token",
				})
				return
			}
		}

		claims, err := utils.ValidateToken(token, utils.GetJWTSecret())
		if err != nil {
			utils.AbortWithError(ctx, fmt.Errorf("validating token: %w", err))
			return
		}
		ctx.Set("claims", claims)
		ctx.Next()
	}
}

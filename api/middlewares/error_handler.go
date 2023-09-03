package middlewares

import (
	"errors"
	"github.com/putto11262002/expense-tracker/api/utils"

	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Assume there will only be one error in the list
// This is because when application encounter an error it will call ctx.Abort immediately
func GlobalErrorHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()

		err := ctx.Errors.Last()
		if err != nil {
			appErr := errors.Unwrap(err)
			log.Printf("ERROR: %v", appErr)

			var dataIntegrityError *utils.DataIntegrityError
			if errors.As(appErr, &dataIntegrityError) {

				ctx.JSON(http.StatusBadRequest, gin.H{"error": dataIntegrityError.Error()})
				return

			}

			var authenticationError *utils.AuthorizationError
			if errors.As(appErr, &authenticationError) {
				ctx.JSON(http.StatusUnauthorized, gin.H{"error": authenticationError.Error()})
				return
			}

			var illegalArgumentError *utils.InvalidArgumentError
			if errors.As(appErr, &illegalArgumentError) {
				ctx.JSON(http.StatusUnauthorized, gin.H{"error": illegalArgumentError.Error()})
				return
			}

			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong, please try again later"})
			return

		}

		// ctx.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("%s %s not found", ctx.Request.Method, ctx.Request.RequestURI)})

	}
}

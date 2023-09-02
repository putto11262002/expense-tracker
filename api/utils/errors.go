package utils

import "github.com/gin-gonic/gin"

// AbortWithError pushes error to gin's error list and abort the request
// Note: don't forget to return after thus function as abort does not terminate the current handler
func AbortWithError(ctx *gin.Context, err error) {
	_ = ctx.Error(err)
	ctx.Abort()
}

// DataIntegrityError errors related to database integrity
type DataIntegrityError struct {
	Message string
}

func (e *DataIntegrityError) Error() string {
	return e.Message
}

// AuthorizationError errors related to authentication
type AuthorizationError struct {
	Message string
}

func (e *AuthorizationError) Error() string {
	return e.Message
}

type InvalidArgumentError struct {
	Message string
}

func (e *InvalidArgumentError) Error() string {
	return e.Message

}

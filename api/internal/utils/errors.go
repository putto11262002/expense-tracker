package utils

import "github.com/gin-gonic/gin"



type DataIntegrityError struct {
	Message string
}

func (e *DataIntegrityError) Error() string {
	return e.Message
}


func AbortWithError(ctx *gin.Context, err error) {
	ctx.Error(err)
	ctx.Abort()
}

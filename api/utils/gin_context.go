package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func GetClaimsFromCtx(ctx *gin.Context) *jwt.StandardClaims {
	val, exist := ctx.Get("claims")
	if !exist {
		return nil
	}
	if claims, ok := val.(*jwt.StandardClaims); !ok {
		return nil
	} else {
		return claims
	}
}

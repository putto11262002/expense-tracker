package middlewares

import (
	"github.com/putto11262002/expense-tracker/api/configs"
	"strings"

	"github.com/gin-gonic/gin"
)

func CORSMiddleware() gin.HandlerFunc {

	allowedOrigins := make(map[string]bool)

	// if in development mode append localhost:3000 by default
	if configs.GetGoEnv() == "development" {
		allowedOrigins["http://localhost:3000"] = true
	}

	// load origin from configs
	originStr, _ := configs.GetStringEnv("ALLOWED_ORIGINS")
	origins := strings.Split(originStr, ",")
	for _, origin := range origins {
		allowedOrigins[origin] = true
	}

	return func(c *gin.Context) {

		if origin := c.Request.Header.Get("Origin"); allowedOrigins[origin] {
			c.Header("Access-Control-Allow-Origin", origin)
		}
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST,HEAD,PATCH, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

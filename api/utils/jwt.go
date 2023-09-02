package utils

import (
	"github.com/putto11262002/expense-tracker/api/configs"
	"github.com/putto11262002/expense-tracker/api/domains"
	"log"
	"time"

	"github.com/golang-jwt/jwt"
)

func GetJWTSecret() string {

	jwtSecret, err := configs.GetStringEnv("JWT_SECRET")
	if err != nil {
		jwtSecret = "secret"
		log.Printf("JWT_SECRET is not configured; using \"secret\" as JWT_SECRET: %v", err)
	}
	return jwtSecret

}

func GenerateJWTToken(user *domains.User, secret string) (string, time.Duration, error) {
	maxAge := time.Minute * 5
	claims := &jwt.StandardClaims{
		Subject:   user.ID.String(),
		ExpiresAt: time.Now().Add(maxAge).Unix(),
		Issuer:    "Expense tracker",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", 0, err
	}
	return ss, maxAge, nil
}

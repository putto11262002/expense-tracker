package utils

import (
	"fmt"
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

func ValidateToken(tokenStr, secret string) (*jwt.StandardClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %s", token.Header["alg"])
		}

		return []byte(secret), nil
	})
	if err != nil {
		return nil, fmt.Errorf("parsing token: %w", err)
	}

	claims, ok := token.Claims.(*jwt.StandardClaims)
	if !ok && !token.Valid {
		return nil, &AuthorizationError{
			Message: "invalid token",
		}
	}

	// check if token has expired
	if err := claims.Valid(); err != nil {
		return nil, &AuthorizationError{
			Message: "invalid token",
		}
	}

	return claims, nil

}

package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/soyaibzihad10/Developer-Assignment/internal/config"
)

type PasswordResetClaims struct {
	UserID  string `json:"user_id"`
	Email   string `json:"email"`
	Purpose string `json:"purpose"`
	jwt.StandardClaims
}

func GeneratePasswordResetToken(userID, email string) (string, error) {
	cfg := config.GetConfig()

	claims := PasswordResetClaims{
		UserID:  userID,
		Email:   email,
		Purpose: "password_reset",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(15 * time.Minute).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.JWT.Secret))
}

func ValidatePasswordResetToken(tokenString string) (string, error) {
	cfg := config.GetConfig()

	token, err := jwt.ParseWithClaims(tokenString, &PasswordResetClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.JWT.Secret), nil
	})

	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok && ve.Errors&jwt.ValidationErrorExpired != 0 {
			return "", fmt.Errorf("password reset token has expired")
		}
		return "", fmt.Errorf("invalid token: %v", err)
	}

	claims, ok := token.Claims.(*PasswordResetClaims)
	if !ok || !token.Valid {
		return "", fmt.Errorf("invalid token claims")
	}

	return claims.UserID, nil
}

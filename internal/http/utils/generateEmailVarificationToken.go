package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/soyaibzihad10/Developer-Assignment/internal/config"
)

type EmailVerificationClaims struct {
	UserID  string `json:"user_id"`
	Email   string `json:"email"`
	Purpose string `json:"purpose"`
	jwt.StandardClaims
}

func GenerateEmailVerificationToken(userID, email string) (string, error) {
	cfg := config.GetConfig()

	claims := EmailVerificationClaims{
		UserID:  userID,
		Email:   email,
		Purpose: "email_verification",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.JWT.Secret))
}

func ValidateEmailVerificationToken(tokenString string) (*EmailVerificationClaims, error) {
	cfg := config.GetConfig()

	token, err := jwt.ParseWithClaims(tokenString, &EmailVerificationClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.JWT.Secret), nil
	})

	if err != nil {
		// Check specifically for expiration error
		if ve, ok := err.(*jwt.ValidationError); ok && ve.Errors&jwt.ValidationErrorExpired != 0 {
			return nil, fmt.Errorf("verification token has expired")
		}
		return nil, err
	}

	claims, ok := token.Claims.(*EmailVerificationClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid verification token")
	}

	// Check if token has expired
	if time.Now().Unix() > claims.ExpiresAt {
		return nil, fmt.Errorf("verification token has expired")
	}

	return claims, nil
}

package utils

import (
	"github.com/tahsin005/affpilot-auth/internal/config"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	cfg := config.GetConfig()
	pass_salt := cfg.PasswordSalt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password + pass_salt), bcrypt.DefaultCost)
	return string(hashedPassword), err
}
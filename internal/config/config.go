package config

import (
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/joho/godotenv"
)

var (
	cfg *Config
	once     sync.Once
)

// GetConfig returns the singleton instance of Config
func GetConfig() *Config {
	once.Do(func() {
		config, err := LoadConfig()
		if err != nil {
			panic(fmt.Sprintf("Failed to load config: %v", err))
		}
		cfg = config
	})
	return cfg
}

// Config holds all configuration for the application
type Config struct {
	App      AppConfig
	Database DatabaseConfig
	JWT      JWTConfig
	Admin    AdminConfig
	Email    EmailConfig
	Server   ServerConfig
}

// AppConfig holds application-specific configuration
type AppConfig struct {
	Environment string
	LogLevel    string
	AppURL      string // Added for password reset and other URLs
}

// DatabaseConfig holds database connection information
type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
	SSLmode  string
}

// JWTConfig holds JWT configuration
type JWTConfig struct {
	Secret string
	Expiry time.Duration
}

// AdminConfig holds system admin information
type AdminConfig struct {
	Username string
	Password string
	Email    string
}

// EmailConfig holds email configuration
type EmailConfig struct {
	VerificationURL string
	From            string
	Host            string
	Port            int
	Username        string
	Password        string
	Secure          bool
	VerificationTTL int
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Port int
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	// Load .env file if it exists
	_ = godotenv.Load()

	// Parse DB port
	dbPort, _ := strconv.Atoi(getEnv("DB_PORT", "5432"))

	// Parse JWT expiry
	jwtExpiry, _ := time.ParseDuration(getEnv("JWT_EXPIRY", "24h"))

	// Parse email port
	emailPort, _ := strconv.Atoi(getEnv("EMAIL_PORT", "587"))

	// Parse email secure
	emailSecure, _ := strconv.ParseBool(getEnv("EMAIL_SECURE", "true"))

	// Parse verification token TTL
	verificationTTL, _ := strconv.Atoi(getEnv("VERIFICATION_TOKEN_TTL", "5"))

	baseURL := getEnv("BASE_URL", "http://localhost:8080")

	// Parse server port
	serverPort, _ := strconv.Atoi(getEnv("SERVER_PORT", "8080"))

	EMAIL_VERIFICATION_URL := fmt.Sprintf("%s/api/v1/auth/verify", baseURL)

	return &Config{
		App: AppConfig{
			Environment: getEnv("APP_ENV", "development"),
			LogLevel:    getEnv("LOG_LEVEL", "debug"),
			AppURL: 	baseURL,
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     dbPort,
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "affpilot"),
			Name:     getEnv("DB_NAME", "affpilot_auth"),
			SSLmode:  getEnv("SSLmode", "affpilot_auth"),
		},
		JWT: JWTConfig{
			Secret: getEnv("JWT_SECRET", "your-secret-key-here"),
			Expiry: jwtExpiry,
		},
		Admin: AdminConfig{
			Username: getEnv("SYSTEM_ADMIN_USERNAME", "admin"),
			Password: getEnv("SYSTEM_ADMIN_PASSWORD", "adminpassword"),
			Email:    getEnv("SYSTEM_ADMIN_EMAIL", "admin@example.com"),
		},
		Email: EmailConfig{
			VerificationURL: EMAIL_VERIFICATION_URL,
			From:            getEnv("EMAIL_FROM", "no-reply@example.com"),
			Host:            getEnv("EMAIL_HOST", "smtp.example.com"),
			Port:            emailPort,
			Username:        getEnv("EMAIL_USERNAME", ""),
			Password:        getEnv("EMAIL_PASSWORD", ""),
			Secure:          emailSecure,
			VerificationTTL: verificationTTL,
		},
		Server: ServerConfig{
			Port: serverPort,
		},
	}, nil
}

// Helper function to get environment variable with a default value
func getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}

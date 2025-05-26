package config

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

var (
	cfg  *Config
	once sync.Once
)

type EmailConfig struct {
	URLBase string
	URLSuffix string
	From string
	Host string
	Port string
	Password string
}

type SystemAdminConfig struct {
	Username string
	Password string
	Email string
}

type Config struct {
	DBUrl string
	Port string
	JWT_SECRET string
	EmailCfg EmailConfig
	SystemAdminCfg SystemAdminConfig
	LoginUrl string
}

type DBConfig struct {
	DBHost string
    DBPort string
    DBUser string
    DBPassword string
    DBName string
}

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}
}

func LoadConfig() (*Config, error) {
	dbConfig := DBConfig {
		DBHost:     os.Getenv("DB_HOST"),
        DBPort:     os.Getenv("DB_PORT"),
        DBUser:     os.Getenv("DB_USER"),
        DBPassword: os.Getenv("DB_PASSWORD"),
        DBName:     os.Getenv("DB_NAME"),
	}

	dbURL := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		dbConfig.DBHost, dbConfig.DBUser, dbConfig.DBPassword, dbConfig.DBName, dbConfig.DBPort,
	)

	port := os.Getenv("SERVER_PORT")

	jwtSecret := os.Getenv("JWT_SECRET")

	emailConfig := EmailConfig {
		URLBase: os.Getenv("EMAIL_VERIFICATION_BASE"),
		URLSuffix: os.Getenv("EMAIL_VERIFICATION_URL_SUFFIX"),
		From: os.Getenv("EMAIL_FROM"),
		Host: os.Getenv("EMAIL_HOST"),
		Port: os.Getenv("EMAIL_PORT"),
		Password: os.Getenv("EMAIL_PASSWORD"),
	}
	
	systemAdminCfg := SystemAdminConfig {
		Username: os.Getenv("SYSTEM_ADMIN_USERNAME"),
		Password: os.Getenv("SYSTEM_ADMIN_PASSWORD"),
		Email: os.Getenv("SYSTEM_ADMIN_EMAIL"),
	}

	loginUrl := os.Getenv("LOGIN_URL")

	if port == "" {
		port = "8080"
	}
	if dbURL == "" {
		log.Fatal("DATABASE_URL not set")
	}

	return &Config{
		DBUrl: dbURL,
		JWT_SECRET: jwtSecret,
		Port:  port,
		EmailCfg: emailConfig,
		SystemAdminCfg: systemAdminCfg,
		LoginUrl: loginUrl,
	}, nil
}

func GetConfig() *Config {
	once.Do(func() {
		var err error
		cfg, err = LoadConfig()
		log.Print(cfg)
		if err != nil {
			log.Fatalf("failed to load config: %v", err)
		}
	})
	return cfg
}

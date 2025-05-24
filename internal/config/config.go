package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the applicationb
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
}

// DatabaseConfig holds database connection information
type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
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
	VerificationURL  string
	PasswordResetURL string
	From             string
	Host             string
	Port             int
	Username         string
	Password         string
	Secure           bool
	VerificationTTL time.Duration

}

// ServerConfig holds server configuration
type ServerConfig struct {
	Port int
}

var (
	cfg  *Config
	once sync.Once
)

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	// Load .env file if it exists
	_ = godotenv.Load()

	// Parse DB port
	dbPort, _ := strconv.Atoi(getEnv("DB_PORT", "5432"))

	// Parse JWT expiry
	jwtExpiryStr := getEnv("JWT_EXPIRY", "24h")
	jwtExpiry, err := time.ParseDuration(jwtExpiryStr)
	if err != nil {
		log.Fatalf("Invalid JWT_EXPIRY value: %v", err)
	}

	// Parse email port
	emailPort, _ := strconv.Atoi(getEnv("EMAIL_PORT", "587"))

	// Parse email secure
	emailSecure, _ := strconv.ParseBool(getEnv("EMAIL_SECURE", "true"))

	// Parse verification token TTL
	// verificationTTL, _ := strconv.Atoi(getEnv("VERIFICATION_TOKEN_TTL", "5"))
	verificationTTLStr := getEnv("VERIFICATION_TOKEN_TTL", "5m")
	verificationTTL, err := time.ParseDuration(verificationTTLStr)
	if err != nil {
		log.Fatalf("Invalid VERIFICATION_TOKEN_TTL value: %v", err)
	}

	// Parse server port
	serverPort, _ := strconv.Atoi(getEnv("SERVER_PORT", "8080"))

	baseURL := getEnv("BASE_URL", "http://localhost")
	port := getEnv("SERVER_PORT", "8080")
	url := fmt.Sprintf("%s:%s", baseURL, port)

	return &Config{
		App: AppConfig{
			Environment: getEnv("APP_ENV", "development"),
			LogLevel:    getEnv("LOG_LEVEL", "debug"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     dbPort,
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			Name:     getEnv("DB_NAME", "affpilot_auth"),
		},
		JWT: JWTConfig{
			Secret: getEnv("JWT_SECRET", "jwtsecretkey"),
			Expiry: jwtExpiry,
		},
		Admin: AdminConfig{
			Username: getEnv("SYSTEM_ADMIN_USERNAME", "admin"),
			Password: getEnv("SYSTEM_ADMIN_PASSWORD", "adminpassword"),
			Email:    getEnv("SYSTEM_ADMIN_EMAIL", "admin@example.com"),
		},
		Email: EmailConfig{
			VerificationURL:  url + "/api/v1/auth/verify",
			PasswordResetURL: url + "/api/v1/auth/reset-password",
			From:             getEnv("EMAIL_FROM", "sagor.sarker0709@gmail.com"),
			Host:             getEnv("EMAIL_HOST", "smtp.gmail.com"),
			Port:             emailPort,
			Username:         getEnv("EMAIL_USERNAME", "smtp"),
			Password:         getEnv("EMAIL_PASSWORD", ""),
			Secure:           emailSecure,
			VerificationTTL:  verificationTTL,
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

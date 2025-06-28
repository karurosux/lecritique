package config

import (
	"fmt"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	App      AppConfig
	Database DatabaseConfig
	Redis    RedisConfig
	JWT      JWTConfig
	Stripe   StripeConfig
	SMTP     *SMTPConfig
	AI       AIConfig
}

type AppConfig struct {
	Name        string
	Env         string
	Port        string
	URL         string
	FrontendURL string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

type JWTConfig struct {
	Secret     string
	Expiration time.Duration
}

type StripeConfig struct {
	SecretKey     string
	WebhookSecret string
}

type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
}

type AIConfig struct {
	Provider string // "openai", "anthropic", or "gemini"
	APIKey   string
	Model    string
}

func Load() (*Config, error) {
	// Load .env file
	_ = godotenv.Load()

	// Set defaults
	viper.SetDefault("APP_NAME", "LeCritique")
	viper.SetDefault("APP_ENV", "development")
	viper.SetDefault("APP_PORT", "8080")
	viper.SetDefault("APP_URL", "http://localhost:8080")
	viper.SetDefault("FRONTEND_URL", "http://localhost:5173")
	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_PORT", "5432")
	viper.SetDefault("DB_SSLMODE", "disable")
	viper.SetDefault("REDIS_HOST", "localhost")
	viper.SetDefault("REDIS_PORT", "6379")
	viper.SetDefault("REDIS_DB", 0)
	viper.SetDefault("JWT_EXPIRATION", "24h")
	viper.SetDefault("AI_PROVIDER", "anthropic")
	viper.SetDefault("AI_MODEL", "claude-3-haiku-20240307")

	// Read from environment
	viper.AutomaticEnv()

	// Parse JWT expiration
	expDuration, err := time.ParseDuration(viper.GetString("JWT_EXPIRATION"))
	if err != nil {
		return nil, fmt.Errorf("invalid JWT_EXPIRATION: %w", err)
	}

	// SMTP config (optional)
	var smtpConfig *SMTPConfig
	if viper.GetString("SMTP_HOST") != "" {
		smtpConfig = &SMTPConfig{
			Host:     viper.GetString("SMTP_HOST"),
			Port:     viper.GetInt("SMTP_PORT"),
			Username: viper.GetString("SMTP_USERNAME"),
			Password: viper.GetString("SMTP_PASSWORD"),
			From:     viper.GetString("SMTP_FROM"),
		}
	}

	config := &Config{
		App: AppConfig{
			Name:        viper.GetString("APP_NAME"),
			Env:         viper.GetString("APP_ENV"),
			Port:        viper.GetString("APP_PORT"),
			URL:         viper.GetString("APP_URL"),
			FrontendURL: viper.GetString("FRONTEND_URL"),
		},
		Database: DatabaseConfig{
			Host:     viper.GetString("DB_HOST"),
			Port:     viper.GetString("DB_PORT"),
			User:     viper.GetString("DB_USER"),
			Password: viper.GetString("DB_PASSWORD"),
			Name:     viper.GetString("DB_NAME"),
			SSLMode:  viper.GetString("DB_SSLMODE"),
		},
		Redis: RedisConfig{
			Host:     viper.GetString("REDIS_HOST"),
			Port:     viper.GetString("REDIS_PORT"),
			Password: viper.GetString("REDIS_PASSWORD"),
			DB:       viper.GetInt("REDIS_DB"),
		},
		JWT: JWTConfig{
			Secret:     viper.GetString("JWT_SECRET"),
			Expiration: expDuration,
		},
		Stripe: StripeConfig{
			SecretKey:     viper.GetString("STRIPE_SECRET_KEY"),
			WebhookSecret: viper.GetString("STRIPE_WEBHOOK_SECRET"),
		},
		SMTP: smtpConfig,
		AI: AIConfig{
			Provider: viper.GetString("AI_PROVIDER"),
			APIKey:   viper.GetString("AI_API_KEY"),
			Model:    viper.GetString("AI_MODEL"),
		},
	}

	return config, nil
}

func (c *Config) DSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Database.Host, c.Database.Port, c.Database.User, c.Database.Password, c.Database.Name, c.Database.SSLMode)
}

// IsSMTPConfigured returns true if SMTP settings are configured
func (c *Config) IsSMTPConfigured() bool {
	return c.SMTP != nil && c.SMTP.Host != "" && c.SMTP.Port > 0
}

// IsDevMode returns true if the app is running in development mode
func (c *Config) IsDevMode() bool {
	return c.App.Env == "development"
}

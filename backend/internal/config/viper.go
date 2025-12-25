package config

import (
	"errors"
	"log"

	"github.com/spf13/viper"
)

func NewViper() *viper.Viper {
	config := viper.New()

	config.SetDefault("APP_NAME", "go-app")
	config.SetDefault("PORT", 8080)
	config.SetDefault("LOG_LEVEL", "info")
	config.SetDefault("DB_HOST", "localhost")
	config.SetDefault("DB_PORT", 5432)
	config.SetDefault("DB_NAME", "go_db")
	config.SetDefault("DB_POOL_IDLE", 10)
	config.SetDefault("DB_POOL_MAX", 100)
	config.SetDefault("DB_POOL_LIFETIME", 300)
	config.SetDefault("JWT_ISSUER", "go-issuer")
	config.SetDefault("JWT_AUDIENCE", "go-audience")
	config.SetDefault("JWT_EXPIRES_MINUTES", 1440)
	config.SetDefault("DROP_TABLE_NAMES", "users,expenses,approvals")
	config.SetDefault("PAYMENT_BASE_URL", "https://1620e98f-7759-431c-a2aa-f449d591150b.mock.pstmn.io")
	config.SetDefault("PAYMENT_TIMEOUT_SECONDS", 10)
	config.SetDefault("PAYMENT_RETRY_COUNT", 3)
	config.SetDefault("PAYMENT_RETRY_DELAY_SECONDS", 2)
	config.SetDefault("PAYMENT_QUEUE_BUFFER", 100)
	config.SetDefault("CORS_ALLOW_ORIGINS", "*")
	config.SetDefault("CORS_ALLOW_CREDENTIALS", false)

	config.SetConfigFile(".env")

	if err := config.ReadInConfig(); err != nil {
		var notFound viper.ConfigFileNotFoundError
		if errors.As(err, &notFound) {
			log.Println("No .env file found in root directory")
		} else {
			log.Printf("Error reading .env file: %v", err)
		}
	} else {
		log.Println("Successfully loaded configuration from .env")
	}

	config.AutomaticEnv()
	return config
}

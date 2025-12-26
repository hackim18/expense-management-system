package config

import (
	"time"

	"github.com/spf13/viper"
)

type paymentConfig struct {
	BaseURL     string
	Timeout     time.Duration
	RetryCount  int
	RetryDelay  time.Duration
	QueueBuffer int
}

func buildPaymentConfig(config *viper.Viper) paymentConfig {
	return paymentConfig{
		BaseURL:     config.GetString("PAYMENT_BASE_URL"),
		Timeout:     time.Duration(config.GetInt("PAYMENT_TIMEOUT_SECONDS")) * time.Second,
		RetryCount:  config.GetInt("PAYMENT_RETRY_COUNT"),
		RetryDelay:  time.Duration(config.GetInt("PAYMENT_RETRY_DELAY_SECONDS")) * time.Second,
		QueueBuffer: config.GetInt("PAYMENT_QUEUE_BUFFER"),
	}
}

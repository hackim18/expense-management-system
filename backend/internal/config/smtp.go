package config

import (
	"go-expense-management-system/internal/integration/email"

	"github.com/spf13/viper"
)

func buildSMTPConfig(config *viper.Viper) email.Config {
	return email.Config{
		Enabled:   config.GetBool("SMTP_ENABLED"),
		Host:      config.GetString("SMTP_HOST"),
		Port:      config.GetInt("SMTP_PORT"),
		Username:  config.GetString("SMTP_USERNAME"),
		Password:  config.GetString("SMTP_PASSWORD"),
		FromEmail: config.GetString("SMTP_FROM_EMAIL"),
		FromName:  config.GetString("SMTP_FROM_NAME"),
	}
}

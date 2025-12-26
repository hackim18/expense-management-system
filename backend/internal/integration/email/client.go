package email

import (
	"context"
	"errors"
	"go-expense-management-system/internal/model"

	"github.com/sirupsen/logrus"
	"gopkg.in/gomail.v2"
)

type Config struct {
	Enabled   bool
	Host      string
	Port      int
	Username  string
	Password  string
	FromEmail string
	FromName  string
}

type Client struct {
	config Config
	log    *logrus.Logger
	dialer *gomail.Dialer
}

func NewClient(config Config, log *logrus.Logger) *Client {
	client := &Client{
		config: config,
		log:    log,
	}
	if config.Enabled && config.Host != "" {
		client.dialer = gomail.NewDialer(config.Host, config.Port, config.Username, config.Password)
	}
	return client
}

func (c *Client) Send(_ context.Context, request model.EmailRequest) error {
	if !c.config.Enabled {
		return nil
	}
	if c.dialer == nil {
		return errors.New("smtp client not configured")
	}
	if len(request.To) == 0 {
		return nil
	}

	fromEmail := c.config.FromEmail
	if fromEmail == "" {
		fromEmail = c.config.Username
	}
	if fromEmail == "" {
		return errors.New("smtp from email is empty")
	}

	message := gomail.NewMessage()
	if c.config.FromName != "" {
		message.SetHeader("From", message.FormatAddress(fromEmail, c.config.FromName))
	} else {
		message.SetHeader("From", fromEmail)
	}
	message.SetHeader("To", request.To...)
	message.SetHeader("Subject", request.Subject)
	message.SetBody("text/plain; charset=UTF-8", request.Body)

	if err := c.dialer.DialAndSend(message); err != nil {
		if c.log != nil {
			c.log.Warnf("Failed to send email: %+v", err)
		}
		return err
	}
	return nil
}

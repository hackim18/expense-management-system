package payment

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"go-expense-management-system/internal/model"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

type Client struct {
	BaseURL    string
	HTTPClient *http.Client
	Log        *logrus.Logger
}

func NewClient(baseURL string, timeout time.Duration, log *logrus.Logger) *Client {
	httpClient := &http.Client{Timeout: timeout}
	return &Client{
		BaseURL:    strings.TrimRight(baseURL, "/"),
		HTTPClient: httpClient,
		Log:        log,
	}
}

func (c *Client) Process(ctx context.Context, request model.PaymentRequest) (*model.PaymentResponse, error) {
	if request.Amount <= 0 || request.ExternalID == "" {
		return nil, fmt.Errorf("invalid payment request")
	}

	payload, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf("%s/v1/payments", c.BaseURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var parsed paymentAPIResponse
	if len(body) > 0 {
		if err := json.Unmarshal(body, &parsed); err != nil {
			return nil, err
		}
	}

	if resp.StatusCode == http.StatusOK {
		return &model.PaymentResponse{
			ID:         parsed.Data.ID,
			ExternalID: parsed.Data.ExternalID,
			Status:     parsed.Data.Status,
		}, nil
	}

	if resp.StatusCode == http.StatusBadRequest {
		message := strings.ToLower(strings.TrimSpace(parsed.Message))
		if strings.Contains(message, "external id already exists") {
			return &model.PaymentResponse{
				ID:         parsed.Data.ID,
				ExternalID: parsed.Data.ExternalID,
				Status:     parsed.Data.Status,
			}, nil
		}
	}

	if c.Log != nil {
		c.Log.Warnf("Payment API error: status=%d body=%s", resp.StatusCode, string(body))
	}
	return nil, fmt.Errorf("payment api error: %s", resp.Status)
}

type paymentAPIResponse struct {
	Data struct {
		ID         string `json:"id"`
		ExternalID string `json:"external_id"`
		Status     string `json:"status"`
	} `json:"data"`
	Message string `json:"message"`
}

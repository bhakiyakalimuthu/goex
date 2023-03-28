package notifier

import (
	"bytes"
	"context"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type Client struct {
	logger     *zap.Logger
	httpClient *http.Client
	url        string
}

func NewHttpClient(logger *zap.Logger) *Client {
	return &Client{logger: logger, httpClient: &http.Client{
		Timeout: 5 * time.Second,
	}}
}

func (c *Client) Notify(ctx context.Context, body string) error {
	req, err := http.NewRequest(http.MethodPost, c.url, bytes.NewBufferString(body))
	if err != nil {
		return err
	}
	_, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	return nil
}

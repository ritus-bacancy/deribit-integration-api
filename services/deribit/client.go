package deribit

import (
	"context"
	"encoding/json"
	"fmt"
	"integration-api/config"
	"net/http"
)

const (
	authEndpoint = "/public/auth"
)

type httpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Client struct {
	httpClient httpClient
	config     *config.Current
}

func NewClient(client httpClient, config *config.Current) *Client {
	return &Client{
		httpClient: client,
		config:     config,
	}
}

func (c *Client) Auth(ctx context.Context) (*Auth, error) {
	params := fmt.Sprintf("client_id=%s&client_secret=%s&grant_type=client_credentials", c.config.ClientID, c.config.ClientSecret)
	url := fmt.Sprintf("%s%s?%s", c.config.BaseURL, authEndpoint, params)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request, %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request, %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed with statuscode %d", resp.StatusCode)
	}

	var auth Auth
	if err := json.NewDecoder(resp.Body).Decode(&auth); err != nil {
		return nil, fmt.Errorf("failed to decode activity %w", err)
	}

	return &auth, nil
}

package deribit

import (
	"context"
	"encoding/json"
	"fmt"
	"integration-api/config"
	"integration-api/request"
	"net/http"
	"strconv"
)

const (
	authEndpoint     = "/public/auth"
	getPriceEndpoint = "/public/get_book_summary_by_currency"
	sellEndpoint     = "/private/sell"
	buyEndpoint      = "/private/buy"
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
	url := fmt.Sprintf("%s%s", c.config.BaseURL, authEndpoint)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request, %w", err)
	}

	q := req.URL.Query()
	q.Add("client_id", c.config.ClientID)
	q.Add("client_secret", c.config.ClientSecret)
	q.Add("grant_type", "client_credentials")
	req.URL.RawQuery = q.Encode()

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request, %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed with statuscode %d", resp.StatusCode)
	}

	var auth Authentication
	if err := json.NewDecoder(resp.Body).Decode(&auth); err != nil {
		return nil, fmt.Errorf("failed to decode auth %w", err)
	}

	return &auth.Result, nil
}

func (c *Client) GetPrice(ctx context.Context, currency string) (*Price, error) {
	url := fmt.Sprintf("%s%s", c.config.BaseURL, getPriceEndpoint)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request, %w", err)
	}

	q := req.URL.Query()
	q.Add("currency", currency)
	q.Add("kind", "future")
	req.URL.RawQuery = q.Encode()

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request, %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed with statuscode %d", resp.StatusCode)
	}

	var price Price
	if err := json.NewDecoder(resp.Body).Decode(&price); err != nil {
		return nil, fmt.Errorf("failed to decode auth %w", err)
	}

	return &price, nil
}

func (c *Client) Buy(ctx context.Context, request request.Buy) (*Buy, error) {
	url := fmt.Sprintf("%s%s", c.config.BaseURL, sellEndpoint)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request, %w", err)
	}

	q := req.URL.Query()
	q.Add("amount", strconv.FormatFloat(request.Amount, 'f', 2, 64))
	q.Add("instrument_name", request.Currency)
	q.Add("label", "market0000234")
	q.Add("type", "market")
	req.URL.RawQuery = q.Encode()

	var bearer = "Bearer " + request.Token
	req.Header.Add("Authorization", bearer)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request, %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed with statuscode %d", resp.StatusCode)
	}

	var buy Buy
	if err := json.NewDecoder(resp.Body).Decode(&buy); err != nil {
		return nil, fmt.Errorf("failed to decode buy %w", err)
	}

	return &buy, nil
}

func (c *Client) Sell(ctx context.Context, request request.Sell) (*Sell, error) {
	url := fmt.Sprintf("%s%s", c.config.BaseURL, buyEndpoint)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request, %w", err)
	}

	q := req.URL.Query()
	q.Add("amount", strconv.FormatFloat(request.Amount, 'f', 2, 64))
	q.Add("instrument_name", request.Currency)
	q.Add("price", strconv.FormatFloat(request.Price, 'f', 2, 64))
	q.Add("trigger", "last_price")
	q.Add("trigger_price", strconv.FormatFloat(request.Price, 'f', 2, 64))
	req.URL.RawQuery = q.Encode()

	var bearer = "Bearer " + request.Token
	req.Header.Add("Authorization", bearer)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request, %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed with statuscode %d", resp.StatusCode)
	}

	var sell Sell
	if err := json.NewDecoder(resp.Body).Decode(&sell); err != nil {
		return nil, fmt.Errorf("failed to decode sell %w", err)
	}

	return &sell, nil
}

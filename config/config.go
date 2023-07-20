package config

import (
	"os"
)

const (
	BaseURL             = "BASE_URL"
	BaseURLDefaultValue = "https://test.deribit.com/api/v2"

	Address             = "ADDRESS"
	AddressDefaultValue = "localhost:8080"

	ClientID = "CLIENT_ID"

	ClientSecret = "CLIENT_SECRET"
)

type Current struct {
	Address      string
	BaseURL      string
	ClientID     string
	ClientSecret string
}

func Load() *Current {
	return &Current{
		BaseURL:      getString(BaseURL, BaseURLDefaultValue),
		Address:      getString(Address, AddressDefaultValue),
		ClientID:     os.Getenv(ClientID),
		ClientSecret: os.Getenv(ClientSecret),
	}
}

func getString(key string, defaultValue string) string {
	val := os.Getenv(key)
	if val == "" {
		return defaultValue
	}
	return val
}

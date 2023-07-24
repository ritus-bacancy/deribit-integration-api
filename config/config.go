package config

import (
	"log"
	"os"
)

const (
	BaseURL             = "BASE_URL"
	BaseURLDefaultValue = "https://test.deribit.com/api/v2"

	Address             = "ADDRESS"
	AddressDefaultValue = ":8080"

	ClientID = "CLIENT_ID"

	ClientSecret = "CLIENT_SECRET"

	Password = "PASSWORD"
)

type Current struct {
	Address      string
	BaseURL      string
	ClientID     string
	ClientSecret string
	Password     string
}

func Load() *Current {
	return &Current{
		BaseURL:      getString(BaseURL, BaseURLDefaultValue),
		Address:      getString(Address, AddressDefaultValue),
		ClientID:     notEmpty(ClientID, os.Getenv(ClientID)),
		ClientSecret: notEmpty(ClientSecret, os.Getenv(ClientSecret)),
		Password:     notEmpty(Password, os.Getenv(Password)),
	}
}

func getString(key string, defaultValue string) string {
	val := os.Getenv(key)
	if val == "" {
		return defaultValue
	}
	return val
}

func notEmpty(key string, val string) string {
	if val == "" {
		log.Fatalf("value is not provided for %s", key)
	}
	return val
}

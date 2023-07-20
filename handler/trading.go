package handler

import (
	"context"
	"integration-api/request"
	"integration-api/services/deribit"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

type service interface {
	Auth(ctx context.Context) (*deribit.Auth, error)
	GetPrice(ctx context.Context, currency string) (*deribit.Price, error)
	Buy(ctx context.Context, req request.Buy) (*deribit.Buy, error)
	Sell(ctx context.Context, req request.Sell) (*deribit.Sell, error)
}

type Trading struct {
	deribitClient service
}

func NewTrading(deribitClient service) *Trading {
	return &Trading{
		deribitClient: deribitClient,
	}
}

func (t *Trading) Auth(c echo.Context) error {
	auth, err := t.deribitClient.Auth(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"access_token":  auth.AccessToken,
		"refresh_token": auth.RefreshToken,
		"expitres_in":   auth.ExpiresIn,
	})
}

func (t *Trading) GetPrice(c echo.Context) error {
	currency := c.QueryParam("currency")
	if currency == "" {
		log.Printf("empty currency in request")
		return c.JSON(http.StatusBadRequest, "empty currency")
	}

	buy, err := t.deribitClient.GetPrice(c.Request().Context(), currency)
	if err != nil {
		log.Printf("error in buy product, %s", err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, buy)
}

func (t *Trading) Buy(c echo.Context) error {
	var req request.Buy
	if err := c.Bind(&req); err != nil {
		log.Printf("failed to parse request, %s", err.Error())
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	req.Token = c.Request().Header.Get("access_token")

	if err := c.Validate(&req); err != nil {
		log.Printf("failed to validate request, %s", err.Error())
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	buy, err := t.deribitClient.Buy(c.Request().Context(), req)
	if err != nil {
		log.Printf("error in deribitClient.Buy, %s", err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, buy)
}

func (t *Trading) Sell(c echo.Context) error {
	var req request.Sell
	if err := c.Bind(&req); err != nil {
		log.Printf("error in sell request, %s", err.Error())
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	req.Token = c.Request().Header.Get("access_token")

	if err := c.Validate(&req); err != nil {
		log.Printf("failed to validate request, %s", err.Error())
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	sell, err := t.deribitClient.Sell(c.Request().Context(), req)
	if err != nil {
		log.Printf("error in sell product, %s", err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, sell)
}

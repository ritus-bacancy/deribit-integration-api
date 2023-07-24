package handler

import (
	"context"
	"fmt"
	"integration-api/config"
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
	config        *config.Current
}

func NewTrading(deribitClient service, config *config.Current) *Trading {
	return &Trading{
		deribitClient: deribitClient,
		config:        config,
	}
}

// Auth godoc
// @summary Authentication endpoint
// @description Authentication endpoint
// @tags auth
// @id Auth
// @accept json
// @produce json
// @Router /auth [POST]
// @response 200 {object} deribit.Auth "OK"
// @Failure  500
func (t *Trading) Auth(c echo.Context) error {
	auth, err := t.deribitClient.Auth(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, auth)
}

// GetPrice godoc
// @summary Get the currency price
// @description Get the currency price
// @tags price
// @id GetPrice
// @Param currency query string  false  "Currency"
// @produce json
// @Router /price [GET]
// @response 200 {object} deribit.Price "OK"
// @Success  200 {object} deribit.Price "OK"
// @Failure  500
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

// Buy godoc
// @summary Buy endpoint
// @description Buy endpoint
// @tags buy
// @id Buy
// @param Buy body request.Buy true "Body"
// @accept json
// @produce json
// @Router /buy [POST]
// @response 200 {object} deribit.Buy "OK"
// @Failure  400
// @Failure  500
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

// Sell godoc
// @summary Sell endpoint
// @description Sell endpoint
// @tags sell
// @id Sell
// @param Sell body request.Sell true "Body"
// @accept json
// @produce json
// @Router /sell [POST]
// @response 200 {object} deribit.Sell "OK"
// @Failure  400
// @Failure  500
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

// Webhook godoc
// @summary Webhook endpoint
// @description Webhook endpoint
// @tags webhook
// @id Webhook
// @param Webhook body request.Webhook true "Body"
// @accept json
// @produce json
// @Router /webhook [POST]
// @response 200 {object} deribit.Sell "OK"
// @response 200 {object} deribit.Buy "OK"
// @Failure  400
// @Failure  500
func (t *Trading) Webhook(c echo.Context) error {
	var req request.Webhook
	if err := c.Bind(&req); err != nil {
		log.Printf("error in sell request, %s", err.Error())
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if req.Password != t.config.Password {
		log.Printf("invalid password ")
		return c.JSON(http.StatusBadRequest, "invalid password ")
	}

	switch req.Operation {
	case "sell":
		auth, err := t.deribitClient.Auth(c.Request().Context())
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		request := request.Sell{
			Token:    auth.AccessToken,
			Amount:   req.Amount,
			Currency: req.Currency,
			Price:    req.Price,
		}
		if err := c.Validate(&request); err != nil {
			log.Printf("failed to validate request, %s", err.Error())
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		sell, err := t.deribitClient.Sell(c.Request().Context(), request)
		if err != nil {
			log.Printf("error in sell product, %s", err.Error())
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, sell)
	case "buy":
		auth, err := t.deribitClient.Auth(c.Request().Context())
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		request := request.Buy{
			Token:    auth.AccessToken,
			Amount:   req.Amount,
			Currency: req.Currency,
		}
		if err := c.Validate(&req); err != nil {
			log.Printf("failed to validate request, %s", err.Error())
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		buy, err := t.deribitClient.Buy(c.Request().Context(), request)
		if err != nil {
			log.Printf("error in deribitClient.Buy, %s", err.Error())
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, buy)
	default:
		log.Printf("operation must be either sell or buy, operation = %s", req.Operation)
		return c.JSON(http.StatusBadRequest, fmt.Errorf("operation must be either sell or buy"))
	}
}

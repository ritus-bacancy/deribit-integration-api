package handler

import (
	"context"
	"integration-api/services/deribit"
	"net/http"

	"github.com/labstack/echo/v4"
)

type service interface {
	Auth(ctx context.Context) (*deribit.Auth, error)
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

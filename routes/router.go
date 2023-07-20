package routes

import (
	"integration-api/config"
	"integration-api/handler"
	"integration-api/services/deribit"
	"net/http"

	"github.com/labstack/echo/v4"
)

type App struct {
	e *echo.Echo

	handlers struct {
		trading *handler.Trading
	}
	services struct {
		deribitService *deribit.Client
	}
}

func (a *App) initialize(config *config.Current) {
	a.initServices(config)
	a.initHandlers()
}

func (a *App) initServices(config *config.Current) {
	a.services.deribitService = deribit.NewClient(http.DefaultClient, config)
}

func (a *App) initHandlers() {
	a.handlers.trading = handler.NewTrading(a.services.deribitService)
}

func (a *App) Set(e *echo.Echo, config *config.Current) {
	a.initialize(config)

	e.POST("/auth", a.handlers.trading.Auth)
}

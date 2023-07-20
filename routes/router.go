package routes

import (
	"integration-api/config"
	"integration-api/handler"
	"integration-api/request"
	"integration-api/services/deribit"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"

	_ "integration-api/docs"

	echoSwagger "github.com/swaggo/echo-swagger"
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
	e.Validator = &request.CustomValidator{Validator: validator.New()}

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.POST("/auth", a.handlers.trading.Auth)
	e.GET("/price", a.handlers.trading.GetPrice)
	e.POST("/buy", a.handlers.trading.Buy)
	e.POST("/sell", a.handlers.trading.Sell)
}

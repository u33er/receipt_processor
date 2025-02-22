package api

import (
	openapiMW "github.com/go-openapi/runtime/middleware"
	"github.com/labstack/echo/v4"
	echoMW "github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
	"net/http"
	"ticket-processor/internal/api/handlers"
	"ticket-processor/internal/api/middlewares"
	"ticket-processor/internal/config"
)

func SetupRouter(log *zap.Logger, cfg *config.Config, h handlers.ReceiptHandler) *echo.Echo {
	e := echo.New()

	e.Use(echoMW.Logger())
	e.Use(echoMW.RequestID())
	e.Use(echoMW.Recover())
	e.Use(echoMW.TimeoutWithConfig(echoMW.TimeoutConfig{
		Timeout: cfg.HTTPServer.Timeout,
		Skipper: echoMW.DefaultSkipper,
	}))
	e.Use(middlewares.ZapLoggerMiddleware(log))

	opts := openapiMW.SwaggerUIOpts{
		SpecURL: "/swagger.json",
		Path:    "swagger",
		Title:   "Ticket Processor API",
	}
	swaggerHandler := openapiMW.SwaggerUI(opts, nil)

	e.GET("/swagger.json", func(c echo.Context) error {
		swagger, err := GetSwagger()
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Error loading Swagger spec").SetInternal(err)
		}

		return c.JSON(http.StatusOK, swagger)
	})

	e.GET("/swagger/*", echo.WrapHandler(swaggerHandler))

	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"message": "Welcome to the Ticket Processor API"})
	})

	e.POST("/receipts/process", h.PostReceiptsProcess)
	e.GET("/receipts/:id/points", h.GetReceiptsIdPoints)

	return e
}

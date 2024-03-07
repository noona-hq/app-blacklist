package server

import (
	"github.com/labstack/echo/v4"
)

func (s Server) NewRouter() *echo.Echo {
	e := echo.New()

	e.GET("/oauth/callback", s.OAuthCallbackHandler)
	e.POST("/webhook", s.WebhookHandler)
	e.GET("/healthz", s.HealthCheckHandler)

	return e
}

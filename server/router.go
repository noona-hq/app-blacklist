package server

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

func (s Server) checkBearerToken(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		auth := c.Request().Header.Get("Authorization")
		if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		}

		// remove "Bearer " prefix
		token := strings.TrimPrefix(auth, "Bearer ")

		if token != s.config.Noona.AppWebhookToken {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		}

		return next(c)
	}
}

func (s Server) NewRouter() *echo.Echo {
	e := echo.New()

	e.GET("/oauth/callback", s.OAuthCallbackHandler)
	e.POST("/webhook", s.WebhookHandler, s.checkBearerToken)
	e.GET("/healthz", s.HealthCheckHandler)

	return e
}

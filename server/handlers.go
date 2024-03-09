package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
	noona "github.com/noona-hq/noona-sdk-go"
)

type SuccessScreenData struct {
	AppStoreURL string
	CompanyName string
}

func (s Server) OAuthCallbackHandler(ctx echo.Context) error {
	data := SuccessScreenData{
		AppStoreURL: s.config.Noona.AppStoreURL,
		CompanyName: "company",
	}

	code := ctx.QueryParam("code")
	if code == "" {
		return ctx.Render(http.StatusOK, "success.html", data)
	}

	user, err := s.services.Core().OnboardUser(code)
	if err != nil {
		s.logger.Errorw("Error onboarding user to app", "error", err)
		return ctx.String(http.StatusInternalServerError, "Something went wrong. Please try again.")
	}

	data.CompanyName = getCompanyNameFromUser(user)

	return ctx.Render(http.StatusOK, "success.html", data)
}

func (s Server) WebhookHandler(ctx echo.Context) error {
	callbackData := noona.CallbackData{}
	if err := ctx.Bind(&callbackData); err != nil {
		s.logger.Errorw("Error binding webhook callback data", "error", err)
		return ctx.String(http.StatusBadRequest, "Bad request")
	}

	if err := s.services.Core().ProcessWebhookCallback(callbackData); err != nil {
		s.logger.Errorw("Error processing webhook callback", "error", err)
		return ctx.String(http.StatusInternalServerError, "Internal server error")
	}

	return ctx.String(http.StatusOK, "WebhookHandler response")
}

func (s Server) HealthCheckHandler(ctx echo.Context) error {
	return ctx.String(http.StatusOK, "OK")
}

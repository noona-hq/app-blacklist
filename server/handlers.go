package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
	noona "github.com/noona-hq/noona-sdk-go"
)

const (
	OpenAction      = "open"
	UninstallAction = "uninstall"
)

type SuccessScreenData struct {
	AppStoreURL string
	CompanyName string
}

func (s Server) OAuthCallbackHandler(ctx echo.Context) error {
	code := ctx.QueryParam("code")
	if code != "" {
		return s.onboardUser(ctx, code)
	}

	IDToken := ctx.QueryParam("id_token")
	if IDToken == "" {
		return ctx.Render(http.StatusOK, "success.html", SuccessScreenData{
			AppStoreURL: s.config.Noona.AppStoreURL,
			CompanyName: "company",
		})
	}

	action := OpenAction
	if ctx.QueryParam("action") != "" {
		action = ctx.QueryParam("action")
	}

	if action == UninstallAction {
		return s.uninstallApp(ctx, IDToken)
	}

	return s.showAppDescription(ctx, IDToken)
}

func (s Server) onboardUser(ctx echo.Context, code string) error {
	_, err := s.services.Core().OnboardUser(code)
	if err != nil {
		s.logger.Errorw("Error onboarding user to app", "error", err)
		return ctx.String(http.StatusInternalServerError, "Something went wrong. Please try again.")
	}

	return ctx.Redirect(http.StatusTemporaryRedirect, s.config.Noona.AppStoreURL)
}

func (s Server) showAppDescription(ctx echo.Context, IDToken string) error {
	user, err := s.services.Core().GetUserFromIDToken(IDToken)
	if err != nil {
		s.logger.Errorw("Error getting user from ID token", "error", err)
		return ctx.String(http.StatusInternalServerError, "Something went wrong. Please try again.")
	}

	data := SuccessScreenData{
		AppStoreURL: s.config.Noona.AppStoreURL,
		CompanyName: getCompanyNameFromUser(user),
	}

	return ctx.Render(http.StatusOK, "success.html", data)
}

func (s Server) uninstallApp(ctx echo.Context, IDToken string) error {
	if err := s.services.Core().UninstallApp(IDToken); err != nil {
		s.logger.Errorw("Error uninstalling app for user", "error", err)
		return ctx.String(http.StatusInternalServerError, "Something went wrong. Please try again.")
	}

	return ctx.String(http.StatusOK, "App uninstalled")
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

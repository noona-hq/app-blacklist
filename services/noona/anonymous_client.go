package noona

import (
	"context"
	"net/http"

	noona "github.com/noona-hq/noona-sdk-go"
	"github.com/pkg/errors"
)

type AnonymousClient struct {
	cfg    Config
	Client *noona.ClientWithResponses
}

func (n AnonymousClient) CodeTokenExchange(code string) (*noona.OAuthToken, error) {
	tokenResponse, err := n.Client.GetOAuthTokenWithResponse(context.Background(), &noona.GetOAuthTokenParams{
		ClientId:     n.cfg.ClientID,
		ClientSecret: n.cfg.ClientSecret,
	}, noona.GetOAuthTokenJSONRequestBody{
		Code:      &code,
		GrantType: noona.AuthorizationCode,
	})
	if err != nil {
		return nil, errors.Wrap(err, "Error in code exchange")
	}

	if tokenResponse.StatusCode() != http.StatusOK {
		return nil, errors.New("Error in code exchange")
	}

	return tokenResponse.JSON200, nil
}

func (n AnonymousClient) RefreshTokenExchange(refreshToken string) (*noona.OAuthToken, error) {
	tokenResponse, err := n.Client.GetOAuthTokenWithResponse(context.Background(), &noona.GetOAuthTokenParams{
		ClientId:     n.cfg.ClientID,
		ClientSecret: n.cfg.ClientSecret,
	}, noona.GetOAuthTokenJSONRequestBody{
		RefreshToken: &refreshToken,
		GrantType:    noona.RefreshToken,
	})
	if err != nil {
		return nil, errors.Wrap(err, "Error in refresh token exchange")
	}

	if tokenResponse.StatusCode() != http.StatusOK {
		return nil, errors.New("Error in refresh token exchange")
	}

	return tokenResponse.JSON200, nil
}

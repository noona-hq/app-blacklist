package noona

import (
	"context"
	"net/http"
	"time"

	"github.com/noona-hq/blacklist/logger"
	"github.com/noona-hq/blacklist/services/store"
	"github.com/noona-hq/blacklist/services/store/entity"
	noona "github.com/noona-hq/noona-sdk-go"
	"github.com/pkg/errors"
)

type Service struct {
	cfg    Config
	logger logger.Logger
	store  store.Store
}

func New(cfg Config, logger logger.Logger, store store.Store) Service {
	return Service{cfg, logger, store}
}

func (s Service) NoAuthNoonaClient() (NoAuthClient, error) {
	client, err := noona.NewClientWithResponses(s.cfg.BaseURL)
	if err != nil {
		return NoAuthClient{}, errors.Wrap(err, "Error creating no auth Noona client")
	}

	return NoAuthClient{Client: client, cfg: s.cfg}, nil
}

func (s Service) AuthNoonaClient(token noona.OAuthToken) (AuthClient, error) {
	if token.AccessToken == nil {
		return AuthClient{}, errors.New("No access token in OAuth token")
	}

	authHeader := noona.WithRequestEditorFn(func(ctx context.Context, req *http.Request) error {
		req.Header.Add("Authorization", "Bearer "+*token.AccessToken)
		return nil
	})

	client, err := noona.NewClientWithResponses(s.cfg.BaseURL, authHeader)
	if err != nil {
		return AuthClient{}, errors.Wrap(err, "Error creating auth Noona client")
	}

	return AuthClient{Client: client, cfg: s.cfg}, nil
}

func (s Service) AuthNoonaClientFromBlacklistUser(blacklistUser entity.User) (AuthClient, error) {
	// Refresh token if it's expired or about to expire
	if blacklistUser.Token.AccessTokenExpiresAt.Before(time.Now().Add(time.Minute * 5)) {
		noAuthClient, err := s.NoAuthNoonaClient()
		if err != nil {
			return AuthClient{}, errors.Wrap(err, "Error getting no auth Noona client")
		}

		token, err := noAuthClient.RefreshTokenExchange(blacklistUser.Token.RefreshToken)
		if err != nil {
			return AuthClient{}, errors.Wrap(err, "Error refreshing token")
		}

		if _, err := s.store.UpdateBlacklistUser(blacklistUser.ID, entity.User{Token: entity.Token{
			AccessToken:          *token.AccessToken,
			AccessTokenExpiresAt: *token.ExpiresAt,
		}}); err != nil {
			s.logger.Errorw("Error updating blacklist user", "error", err)
		}

		return s.AuthNoonaClient(*token)
	}

	return s.AuthNoonaClient(noona.OAuthToken{AccessToken: &blacklistUser.Token.AccessToken})
}

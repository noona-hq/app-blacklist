package noona

import (
	"github.com/noona-hq/app-blacklist/logger"
	"github.com/noona-hq/app-blacklist/store"
	noona "github.com/noona-hq/noona-sdk-go"
	"github.com/pkg/errors"
)

type Service struct {
	cfg    Config
	logger logger.Logger
	store  store.Store
}

func New(cfg Config, logger logger.Logger, store store.Store) Service {
	return Service{cfg: cfg, logger: logger, store: store}
}

func (s Service) AnonymousClient() (AnonymousClient, error) {
	client, err := noona.NewAnonymous(noona.ClientOptions{
		BaseURL: s.cfg.BaseURL,
	})
	if err != nil {
		return AnonymousClient{}, errors.Wrap(err, "Error creating anonymous Noona client")
	}

	return AnonymousClient{Client: client, cfg: s.cfg}, nil
}

func (s Service) Client(token noona.OAuthToken) (Client, error) {
	if token.AccessToken == nil {
		return Client{}, errors.New("No access token in OAuth token")
	}

	client, err := noona.New(*token.AccessToken, noona.ClientOptions{
		BaseURL: s.cfg.BaseURL,
	})
	if err != nil {
		return Client{}, errors.Wrap(err, "Error creating auth Noona client")
	}

	return Client{Client: client, cfg: s.cfg}, nil
}

func (s Service) ClientID() string {
	return s.cfg.ClientID
}

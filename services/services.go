package services

import (
	"github.com/noona-hq/blacklist/logger"
	"github.com/noona-hq/blacklist/services/core"
	"github.com/noona-hq/blacklist/services/noona"
	"github.com/noona-hq/blacklist/store"
	"github.com/pkg/errors"
)

type Services struct {
	logger logger.Logger
	core   core.Service
	noona  noona.Service
}

func New(noonaCfg noona.Config, logger logger.Logger, store store.Store) (Services, error) {
	noonaService := noona.New(noonaCfg, logger, store)
	coreService, err := core.New(logger, noonaService, store)
	if err != nil {
		return Services{}, errors.Wrap(err, "error creating core service")
	}

	return Services{
		logger: logger,
		core:   coreService,
		noona:  noonaService,
	}, nil
}

func (s *Services) Noona() noona.Service {
	return s.noona
}

func (s *Services) Core() core.Service {
	return s.core
}

package services

import (
	"github.com/noona-hq/blacklist/logger"
	"github.com/noona-hq/blacklist/services/core"
	"github.com/noona-hq/blacklist/services/noona"
	"github.com/noona-hq/blacklist/services/store"
)

type Services struct {
	logger logger.Logger
	core   core.Service
	noona  noona.Service
}

func New(noonaCfg noona.Config, logger logger.Logger, store store.Store) Services {
	noonaService := noona.New(noonaCfg, logger)

	return Services{
		core:  core.New(logger, noonaService, store),
		noona: noonaService,
	}
}

func (s *Services) Noona() noona.Service {
	return s.noona
}

func (s *Services) Core() core.Service {
	return s.core
}

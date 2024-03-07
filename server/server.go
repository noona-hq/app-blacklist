package server

import (
	"log"
	"net/http"

	"github.com/noona-hq/blacklist/db"
	"github.com/noona-hq/blacklist/logger"
	"github.com/noona-hq/blacklist/server/templates"
	"github.com/noona-hq/blacklist/services"
	"github.com/noona-hq/blacklist/services/store"
	"github.com/noona-hq/blacklist/services/store/memory"
	"github.com/noona-hq/blacklist/services/store/mongodb"
	"github.com/pkg/errors"
)

type Server struct {
	config   Config
	logger   logger.Logger
	services services.Services
}

func New(config Config, logger logger.Logger) (Server, error) {
	server := Server{
		config: config,
		logger: logger,
	}

	var store store.Store
	var err error
	switch config.Store {
	case "mongodb":
		store, err = server.MongoStore()
		if err != nil {
			return Server{}, errors.Wrap(err, "unable to create mongodb store")
		}
	case "memory":
		store = server.MemoryStore()
	default:
		store, err = server.MongoStore()
		if err != nil {
			return Server{}, errors.Wrap(err, "unable to create mongodb store")
		}
	}

	server.services = services.New(config.Noona, logger, store)

	return server, nil
}

func (s *Server) Serve() error {
	router := s.NewRouter()
	router.Renderer = templates.NewRenderer(s.logger)

	log.Println("Starting Blacklist server...")
	return http.ListenAndServe(":8080", router)
}

func (s *Server) MongoStore() (store.Store, error) {
	database, err := db.New(s.config.DB, s.logger)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create database")
	}

	return mongodb.NewStore(*database), nil
}

func (s *Server) MemoryStore() store.Store {
	return memory.NewStore()
}

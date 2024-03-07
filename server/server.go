package server

import (
	"log"
	"net/http"

	"github.com/noona-hq/blacklist/db"
	"github.com/noona-hq/blacklist/logger"
	"github.com/noona-hq/blacklist/server/templates"
	"github.com/noona-hq/blacklist/services"
	"github.com/noona-hq/blacklist/services/store/mongodb"
	"github.com/pkg/errors"
)

type Server struct {
	config   Config
	logger   logger.Logger
	services services.Services
}

func New(config Config, logger logger.Logger) (Server, error) {
	database, err := db.New(config.DB, logger)
	if err != nil {
		return Server{}, errors.Wrap(err, "unable to create database")
	}

	store := mongodb.NewStore(*database)

	return Server{
		config:   config,
		logger:   logger,
		services: services.New(config.Noona, logger, store),
	}, nil
}

func (s *Server) Serve() error {
	router := s.NewRouter()
	router.Renderer = templates.NewRenderer(s.logger)

	log.Println("Starting Blacklist server...")
	return http.ListenAndServe(":8080", router)
}

package server

import (
	"github.com/noona-hq/app-blacklist/db"
	"github.com/noona-hq/app-blacklist/logger"
	"github.com/noona-hq/app-blacklist/services/noona"
)

type Config struct {
	Noona  noona.Config
	Logger logger.Config
	DB     db.Config
	// Store can either be memory or mongodb
	Store string `default:"mongodb"`
}

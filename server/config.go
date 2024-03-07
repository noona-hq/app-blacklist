package server

import (
	"github.com/noona-hq/blacklist/db"
	"github.com/noona-hq/blacklist/logger"
	"github.com/noona-hq/blacklist/services/noona"
)

type Config struct {
	Noona  noona.Config
	Logger logger.Config
	DB     db.Config
}

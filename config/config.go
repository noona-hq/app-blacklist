package config

import (
	"github.com/kelseyhightower/envconfig"
)

func Process(cfg interface{}) error {
	return envconfig.Process("", cfg)
}

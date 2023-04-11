package config

import (
	"flag"
	"github.com/caarlos0/env"
)

type Config struct {
}

var cfg Config

func init() {
}

func New() (Config, error) {
	if err := env.Parse(&cfg); err != nil {
		return Config{}, err
	}
	flag.Parse()
	return cfg, nil
}

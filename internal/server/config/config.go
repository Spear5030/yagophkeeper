package config

import (
	"flag"
	"github.com/caarlos0/env"
)

type Config struct {
	FileStorage string `env:"GK_SERVER_FILE" envDefault:"gkdata.pbb"`
	Secret      string `env:"GK_SERVER_SECRET" envDefault:"V3ry$trongK3y"`
	Port        string `env:"GK_SERVER_PORT" envDefault:"22345"`
	ServerCert  string `env:"GK_SERVER_CERT" envDefault:"cert/server-cert.pem"`
	ServerKey   string `env:"GK_SERVER_KEY" envDefault:"cert/server-key.pem"`
}

var cfg Config

func New() (Config, error) {
	if err := env.Parse(&cfg); err != nil {
		return Config{}, err
	}
	flag.Parse()
	return cfg, nil
}

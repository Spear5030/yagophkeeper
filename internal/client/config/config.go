package config

import "github.com/caarlos0/env"

type Config struct {
	FileStorage string `env:"GK_CLIENT_FILE" envDefault:"user.dat"`
	Addr        string `env:"GK_SERVER_ADDR" envDefault:":22345"`
	Cert        string `env:"GK_CLIENT_CERT" envDefault:"cert/ca-cert.pem"`
}

var cfg Config

func New() (Config, error) {
	if err := env.Parse(&cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

package config

import "flag"
import "github.com/caarlos0/env"

type Config struct {
	Addr string `env:"YAGOPHKEEPER"`
	Cert string `env:"YAGOPHKEEPER_CERT"`
}

var cfg Config

func init() {
	flag.StringVar(&cfg.Addr, "a", cfg.Addr, "Server address")
	flag.StringVar(&cfg.Cert, "c", cfg.Cert, "Cert path")
}

func New() (Config, error) {
	if err := env.Parse(&cfg); err != nil {
		return Config{}, err
	}
	flag.Parse()
	return cfg, nil
}

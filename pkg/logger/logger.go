// Package logger отвечает за создание zap.Logger
package logger

import "go.uber.org/zap"

// New makes new zap.Logger by debug level.
func New(debug bool) (*zap.Logger, error) {
	if debug {
		return zap.NewDevelopment()
	}
	return zap.NewProduction()
}

func NewForClient(debug bool) (*zap.Logger, error) {
	if debug {
		cfg := zap.NewProductionConfig()
		cfg.OutputPaths = []string{
			"yagophkeeper.log", "stderr",
		}
		return cfg.Build()
	}
	return zap.NewProduction()
}

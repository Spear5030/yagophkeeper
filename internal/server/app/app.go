package app

import (
	"github.com/Spear5030/yagophkeeper/internal/client/config"
	"github.com/Spear5030/yagophkeeper/pkg/logger"
	"go.uber.org/zap"
)

type App struct {
	logger *zap.Logger
}

func New(cfg config.Config) (*App, error) {
	lg, err := logger.New(true)
	if err != nil {
		return nil, err
	}
	return &App{logger: lg}, nil
}

func (app *App) Run() error {
	return nil
}

package app

import (
	"github.com/Spear5030/yagophkeeper/internal/server"
	"github.com/Spear5030/yagophkeeper/internal/server/config"
	"github.com/Spear5030/yagophkeeper/internal/server/storage"
	"github.com/Spear5030/yagophkeeper/internal/server/usecase"
	"github.com/Spear5030/yagophkeeper/pkg/logger"
	"go.uber.org/zap"
)

type App struct {
	GRPCServer *server.YaGophKeeperServer
	logger     *zap.Logger
}

func New(cfg config.Config) (*App, error) {
	lg, err := logger.New(true)
	if err != nil {
		return nil, err
	}
	s := storage.New("test.ppb", lg)
	uc := usecase.New(s, lg)
	srv := server.New(uc, lg, "12345") //todo cfg
	return &App{GRPCServer: srv, logger: lg}, nil
}

func (app *App) Run() error {
	return app.GRPCServer.Start()

}

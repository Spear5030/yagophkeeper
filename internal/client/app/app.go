package app

import (
	"github.com/Spear5030/yagophkeeper/internal/client/cli"
	"github.com/Spear5030/yagophkeeper/internal/client/config"
	"github.com/Spear5030/yagophkeeper/internal/client/storage"
	"github.com/Spear5030/yagophkeeper/internal/client/usecase"
	"github.com/Spear5030/yagophkeeper/pkg/logger"
	"go.uber.org/zap"
)

type App struct {
	logger *zap.Logger
	client *cli.CLI
}

func New(cfg config.Config) (*App, error) {
	lg, err := logger.New(true)
	if err != nil {
		return nil, err
	}
	repo, err := storage.New("user.dat", lg)
	useCase := usecase.New(lg, repo)

	client := cli.New(lg, useCase)

	return &App{
			logger: lg,
			client: client},
		nil
}

func (app *App) Run() error {
	app.client.Execute()
	return nil
}

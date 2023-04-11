package app

import (
	"github.com/Spear5030/yagophkeeper/internal/client/cli"
	"github.com/Spear5030/yagophkeeper/internal/client/config"
	"github.com/Spear5030/yagophkeeper/internal/client/grpcclient"
	"github.com/Spear5030/yagophkeeper/internal/client/storage"
	"github.com/Spear5030/yagophkeeper/internal/client/usecase"
	"github.com/Spear5030/yagophkeeper/pkg/logger"
	"go.uber.org/zap"
)

type App struct {
	logger *zap.Logger
	cli    *cli.CLI
	grpc   *grpcclient.Client
}

func New(cfg config.Config) (*App, error) {
	lg, err := logger.New(true)
	if err != nil {
		return nil, err
	}
	repo, err := storage.New("user.dat", lg)
	grpcl := grpcclient.New("localhost:12345", repo.Token) //todo cfg
	useCase := usecase.New(repo, grpcl, lg)
	cliclient := cli.New(lg, useCase)

	return &App{
			logger: lg,
			cli:    cliclient},
		nil
}

func (app *App) Run() error {
	app.cli.Execute()
	return nil
}

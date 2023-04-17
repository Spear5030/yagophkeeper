package app

import (
	"github.com/Spear5030/yagophkeeper/internal/client/cli"
	"github.com/Spear5030/yagophkeeper/internal/client/config"
	"github.com/Spear5030/yagophkeeper/internal/client/grpcclient"
	"github.com/Spear5030/yagophkeeper/internal/client/storage"
	"github.com/Spear5030/yagophkeeper/internal/client/usecase"
	"github.com/Spear5030/yagophkeeper/pkg/logger"
	"go.uber.org/zap"
	"log"
)

type App struct {
	logger *zap.Logger
	cli    *cli.CLI
}

func New(cfg config.Config) (*App, error) {
	lg, err := logger.New(true)
	if err != nil {
		return nil, err
	}
	repo, err := storage.New(cfg.FileStorage, cfg.MasterPass, lg)
	if err != nil {
		log.Fatal(err)
	}
	grpcl := grpcclient.New(cfg.Addr, cfg.Cert, repo.GetToken())
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

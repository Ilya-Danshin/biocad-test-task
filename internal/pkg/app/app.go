package app

import (
	"context"

	"test_task/internal/app/config"
	"test_task/internal/app/database"
	"test_task/internal/app/service"
)

type App struct {
	cfg *config.Config
	db  *database.Postgres
	s   *service.Service
}

func New() (*App, error) {
	a := &App{}
	var err error

	a.cfg, err = config.New()
	if err != nil {
		return nil, err
	}

	a.db, err = database.New(&a.cfg.Database, context.Background())
	if err != nil {
		return nil, err
	}

	a.s, err = service.New()
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) Run() error {
	return nil
}

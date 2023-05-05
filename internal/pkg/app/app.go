package app

import (
	"context"
	"log"
	"test_task/internal/app/http_service"

	"test_task/internal/app/config"
	"test_task/internal/app/database"
	"test_task/internal/app/directory"
	"test_task/internal/app/grpc_service"
	"test_task/internal/app/parser"
)

type App struct {
	cfg  *config.Config
	db   *database.Postgres
	dir  *directory.FilesDirectory
	grpc *grps_service.Service
	http *http_service.Service
	par  *parser.Parser

	errors chan error
}

func New() (*App, error) {
	a := &App{}
	var err error

	ctx := context.Background()

	a.cfg, err = config.New()
	if err != nil {
		return nil, err
	}

	a.db, err = database.New(a.cfg, ctx)
	if err != nil {
		return nil, err
	}

	queue := make(chan string, a.cfg.App.QueueMaxSize)
	a.errors = make(chan error)

	a.dir, err = directory.New(ctx, a.cfg, queue, a.db, a.errors)
	if err != nil {
		return nil, err
	}

	a.par, err = parser.New(a.cfg, queue, a.errors, a.db)
	if err != nil {
		return nil, err
	}

	a.grpc, err = grps_service.New(a.cfg, a.db, a.errors)
	if err != nil {
		return nil, err
	}

	a.http, err = http_service.New(a.cfg, a.db, a.errors)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) Run() error {
	ctx := context.Background()

	go a.dir.Run(ctx)
	go a.par.Run(ctx)
	go a.grpc.Run()
	go a.http.Run()

	for {
		log.Print(<-a.errors)
	}

	return nil
}

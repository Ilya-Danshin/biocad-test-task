package app

import (
	"context"
	"log"

	"test_task/internal/app/config"
	"test_task/internal/app/database"
	"test_task/internal/app/directory"
	"test_task/internal/app/parser"
	"test_task/internal/app/service"
)

type App struct {
	cfg *config.Config
	db  *database.Postgres
	dir *directory.FilesDirectory
	s   *service.Service
	par *parser.Parser

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

	a.db, err = database.New(&a.cfg.Database, ctx)
	if err != nil {
		return nil, err
	}

	queue := make(chan string, 1024)
	a.errors = make(chan error)

	a.dir, err = directory.New(ctx, a.cfg.FilesDirectory, queue, a.db, a.errors)
	if err != nil {
		return nil, err
	}

	a.par, err = parser.New(a.cfg.Parser, queue, a.errors, a.db)
	if err != nil {
		return nil, err
	}

	a.s, err = service.New(a.cfg.Service, a.db, a.errors)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) Run() error {
	ctx := context.Background()

	go a.dir.Run(ctx)
	go a.par.Run(ctx)
	go a.s.Run()

	for {
		log.Print(<-a.errors)
	}

	return nil
}

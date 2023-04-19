package app

import (
	"context"
	"log"

	"test_task/internal/app/config"
	"test_task/internal/app/database"
	"test_task/internal/app/directory"
	"test_task/internal/app/service"
)

type App struct {
	cfg *config.Config
	db  *database.Postgres
	dir *directory.FilesDirectory
	s   *service.Service

	errors chan error
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

	queue := make(chan string, 1024)
	a.errors = make(chan error)

	a.dir, err = directory.New(a.cfg.FilesDirectory, queue, a.db, a.errors)
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

	go a.dir.Run()

	for {
		log.Print(<-a.errors)
	}

	return nil
}

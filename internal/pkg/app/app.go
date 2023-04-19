package app

import "test_task/internal/app/service"

type App struct {
	s *service.Service
}

func New() (*App, error) {
	a := &App{}
	var err error

	a.s, err = service.New()
	if err != nil {
		return nil, err
	}
	
	return a, nil
}

func (a *App) Run() error {
	return nil
}

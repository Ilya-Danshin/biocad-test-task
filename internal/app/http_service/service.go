package http_service

import (
	"context"
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"

	"test_task/internal/app/config"
	"test_task/internal/app/database"
)

type Service struct {
	pageSize int
	db       database.IDatabase
	port     int

	e *echo.Echo

	errChan chan error
}

func New(cfg *config.Config, db database.IDatabase, errChan chan error) (*Service, error) {
	s := &Service{}

	s.pageSize = cfg.HttpService.PageSize
	s.port = cfg.HttpService.Port
	s.db = db
	s.errChan = errChan

	s.e = s.initEchoService()

	return s, nil
}

func (s *Service) initEchoService() *echo.Echo {
	// Echo instance
	e := echo.New()

	e.POST("/", s.GetData)

	return e
}

func (s *Service) GetData(ctx echo.Context) error {

	guid, err := uuid.FromString(ctx.QueryParam("guid"))
	if err != nil {
		internalErrorHandler(ctx, fmt.Sprintf("guid convert from query param error: %s", err.Error()))
		return err
	}

	page, err := strconv.Atoi(ctx.QueryParam("page"))
	if err != nil {
		internalErrorHandler(ctx, fmt.Sprintf("page convert from query param error: %s", err.Error()))
		return err
	}

	limit, err := strconv.Atoi(ctx.QueryParam("limit"))
	if err != nil {
		internalErrorHandler(ctx, fmt.Sprintf("limit convert from query param error: %s", err.Error()))
		return err
	}

	data, err := s.db.GetDataAPI(context.Background(), guid, page*s.pageSize, limit)
	if err != nil {
		internalErrorHandler(ctx, fmt.Sprintf("get data error: %s", err.Error()))
		return err
	}

	err = ctx.JSON(http.StatusOK, data)
	if err != nil {
		return err
	}

	return nil
}

func internalErrorHandler(ctx echo.Context, err string) {
	ctx.String(http.StatusInternalServerError, err)
}

func (s *Service) Run() {
	err := s.e.Start(":" + strconv.Itoa(s.port))
	if err != nil {
		s.errChan <- err
		return
	}

	return
}

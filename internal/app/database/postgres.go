package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"

	"test_task/internal/app/config"
)

type Postgres struct {
	conn *pgx.Conn
}

func New(cfg *config.DB, ctx context.Context) (*Postgres, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		cfg.Host, cfg.User, cfg.Password, cfg.DatabaseName, cfg.Port)

	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		return nil, err
	}

	db := Postgres{}
	db.conn = conn

	return &db, nil
}

func (db *Postgres) Add() {

}

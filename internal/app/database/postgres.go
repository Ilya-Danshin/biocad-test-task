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

func (db *Postgres) AddProcessedFile(ctx context.Context, filename string) error {
	rows, err := db.conn.Query(ctx,
		`INSERT INTO files VALUES ($1)`, filename)
	defer rows.Close()
	if err != nil {
		return err
	}

	return nil
}

func (db *Postgres) GetProcessedFiles(ctx context.Context) ([]string, error) {
	rows, err := db.conn.Query(ctx,
		`SELECT * FROM files;`)
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	var files []string

	for rows.Next() {
		var file string
		err = rows.Scan(&file)
		if err != nil {
			return nil, err
		}

		files = append(files, file)
	}

	return files, nil
}

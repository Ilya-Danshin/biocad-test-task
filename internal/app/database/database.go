package database

import "context"

type IDatabase interface {
	AddProcessedFile(ctx context.Context, filename string) error

	GetProcessedFiles(ctx context.Context) ([]string, error)
}

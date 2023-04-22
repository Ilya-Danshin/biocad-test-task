package outfile

import (
	"context"

	"test_task/internal/app/database"
)

type IOutFile interface {
	WriteData(ctx context.Context, records []database.Record) error
}

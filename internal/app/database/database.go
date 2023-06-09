package database

import (
	"context"

	"github.com/gofrs/uuid"
)

type Record struct {
	N         int
	MQTT      []byte
	InvId     string
	UnitGuid  uuid.UUID
	MsgId     string
	Text      string
	Context   []byte
	Class     string
	Level     int
	Area      string
	Addr      string
	Block     string
	Type      string
	Bit       int
	InvertBit int
}

type IDatabase interface {
	AddProcessedFile(ctx context.Context, filename string) error

	GetProcessedFiles(ctx context.Context) ([]string, error)

	AddDataRow(ctx context.Context, data []Record) error
	GetRecordsByGuid(ctx context.Context, guid uuid.UUID) ([]Record, error)

	GetDataAPI(ctx context.Context, guid uuid.UUID, offset int32, limit int32) ([]Record, error)
}

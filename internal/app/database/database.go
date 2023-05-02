package database

import (
	"context"

	"github.com/balibuild/winio/pkg/guid"
)

type Record struct {
	N         int
	MQTT      []byte
	InvId     string
	UnitGuid  guid.GUID
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
	GetRecordsByGuid(ctx context.Context, guid guid.GUID) ([]Record, error)

	GetDataAPI(ctx context.Context, guid guid.GUID, offset int32, limit int32) ([]Record, error)
}

package database

import (
	"context"
	"fmt"
	"github.com/gofrs/uuid"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"

	"test_task/internal/app/config"
)

type Postgres struct {
	conn *pgxpool.Pool
}

func New(cfg *config.Config, ctx context.Context) (*Postgres, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		cfg.Database.Host, cfg.Database.User, cfg.Database.Password, cfg.Database.DatabaseName, cfg.Database.Port)

	conn, err := pgxpool.Connect(ctx, dsn)
	if err != nil {
		return nil, err
	}

	db := Postgres{}
	db.conn = conn

	return &db, nil
}

func (db *Postgres) AddProcessedFile(ctx context.Context, filename string) error {
	rows, err := db.conn.Query(ctx,
		`INSERT INTO files VALUES ($1);`, filename)
	if err != nil {
		return err
	}
	defer rows.Close()

	return nil
}

func (db *Postgres) GetProcessedFiles(ctx context.Context) ([]string, error) {
	rows, err := db.conn.Query(ctx,
		`SELECT file FROM files;`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

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

const batchQueue string = `INSERT INTO data
							VALUES ($1, $2, $3, $4, $5, 
							        $6, $7, $8, $9, $10, 
							        $11, $12, $13, $14, $15);`

func (db *Postgres) AddDataRow(ctx context.Context, data []Record) error {
	tx, err := db.conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
	}()

	batch := &pgx.Batch{}

	for _, row := range data {
		batch.Queue(batchQueue,
			row.N, row.MQTT, row.InvId, row.UnitGuid, row.MsgId, row.Text, row.Context, row.Class,
			row.Level, row.Area, row.Addr, row.Block, row.Type, row.Bit, row.InvertBit)
	}

	br := tx.SendBatch(ctx, batch)

	for _, _ = range data {
		_, err = br.Exec()
		if err != nil {
			return err
		}
	}

	err = br.Close()
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)

	return nil
}

const recordsByGuidQueue string = `SELECT n, mqtt, invid, unit_guid, msg_id, 
       										text, context, class, level, area, 
       										addr, block, type, bit, invert_bit
									FROM data WHERE (unit_guid=$1);`

func (db *Postgres) GetRecordsByGuid(ctx context.Context, guid uuid.UUID) ([]Record, error) {
	rows, err := db.conn.Query(ctx, recordsByGuidQueue, guid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var allRecords []Record
	for rows.Next() {
		var oneRecord Record
		err = rows.Scan(
			&oneRecord.N,
			&oneRecord.MQTT,
			&oneRecord.InvId,
			&oneRecord.UnitGuid,
			&oneRecord.MsgId,
			&oneRecord.Text,
			&oneRecord.Context,
			&oneRecord.Class,
			&oneRecord.Level,
			&oneRecord.Area,
			&oneRecord.Addr,
			&oneRecord.Block,
			&oneRecord.Type,
			&oneRecord.Bit,
			&oneRecord.InvertBit,
		)
		if err != nil {
			return nil, err
		}

		allRecords = append(allRecords, oneRecord)
	}

	return allRecords, nil
}

const recordsByGuidWithOffsetLimitQueue string = `SELECT n, mqtt, invid, unit_guid, msg_id, 
       												text, context, class, level, area, 
       												addr, block, type, bit, invert_bit
												FROM data WHERE unit_guid=$1 LIMIT $2 OFFSET $3;`

func (db *Postgres) GetDataAPI(ctx context.Context, guid uuid.UUID, offset int32, limit int32) ([]Record, error) {
	rows, err := db.conn.Query(ctx, recordsByGuidWithOffsetLimitQueue, guid, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var allRecords []Record
	for rows.Next() {
		var oneRecord Record
		err = rows.Scan(
			&oneRecord.N,
			&oneRecord.MQTT,
			&oneRecord.InvId,
			&oneRecord.UnitGuid,
			&oneRecord.MsgId,
			&oneRecord.Text,
			&oneRecord.Context,
			&oneRecord.Class,
			&oneRecord.Level,
			&oneRecord.Area,
			&oneRecord.Addr,
			&oneRecord.Block,
			&oneRecord.Type,
			&oneRecord.Bit,
			&oneRecord.InvertBit,
		)
		if err != nil {
			return nil, err
		}

		allRecords = append(allRecords, oneRecord)
	}

	return allRecords, nil
}

package parser

import (
	"context"
	"encoding/csv"
	"os"
	"strconv"
	"strings"

	"github.com/gofrs/uuid"
	"github.com/pkg/errors"

	"test_task/internal/app/config"
	"test_task/internal/app/database"
	"test_task/internal/app/outfile"
)

type Parser struct {
	queue       chan string
	db          database.IDatabase
	outFilesDir string
	pdfApiKey   string
	outFile     outfile.IOutFile

	errChan chan error
}

func New(cfg *config.Config, queue chan string, errChan chan error, db database.IDatabase) (*Parser, error) {
	par := Parser{}

	par.queue = queue
	par.errChan = errChan
	par.db = db
	par.outFilesDir = cfg.Parser.OutFilesDirectory
	par.pdfApiKey = cfg.Parser.PdfApiKey

	var err error
	par.outFile, err = outfile.New(cfg, par.db)
	if err != nil {
		return nil, err
	}

	return &par, nil
}

func (p *Parser) Run(ctx context.Context) {
	for {
		file := <-p.queue
		tsvData, err := p.readTSVFile(file)
		if err != nil {
			p.errChan <- errors.Errorf("read tsv file error: %e", err)
			continue
		}

		records, err := p.parseTSV(tsvData)
		if err != nil {
			p.errChan <- errors.Errorf("parse tsv file error: %e", err)
		}

		err = p.db.AddDataRow(ctx, records)
		if err != nil {
			p.errChan <- errors.Errorf("add data to database error: %e", err)
		}

		err = p.WriteDataToFile(ctx, records)
		if err != nil {
			p.errChan <- errors.Errorf("write to out file error: %e", err)
		}
	}
}

func (p *Parser) readTSVFile(path string) (*[][]string, error) {
	tsvFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	r := csv.NewReader(tsvFile)
	r.Comma = '\t' // Use tab-delimited instead of comma

	tsvData, err := r.ReadAll()
	if err != nil {
		return nil, err
	}

	err = tsvFile.Close()
	if err != nil {
		return nil, err
	}

	return &tsvData, nil
}

func (p *Parser) parseTSV(tsvData *[][]string) ([]database.Record, error) {
	var allRecords []database.Record
	var err error

	for _, row := range *tsvData {
		var oneRecord database.Record

		oneRecord.N, err = readInt(row[0])
		if err != nil {
			p.errChan <- err
			continue
		}

		oneRecord.MQTT, _ = readBytes(row[1])
		oneRecord.InvId, _ = readString(row[2])

		oneRecord.UnitGuid, err = uuid.FromString(row[3])
		if err != nil {
			p.errChan <- err
			continue
		}

		oneRecord.MsgId, _ = readString(row[4])
		oneRecord.Text, _ = readString(row[5])
		oneRecord.Context, _ = readBytes(row[6])
		oneRecord.Class, _ = readString(row[7])

		oneRecord.Level, err = readInt(row[8])
		if err != nil {
			p.errChan <- err
			continue
		}

		oneRecord.Area, _ = readString(row[9])
		oneRecord.Addr, _ = readString(row[10])
		oneRecord.Block, _ = readString(row[11])
		oneRecord.Type, _ = readString(row[12])

		oneRecord.Bit, err = readInt(row[13])
		if err != nil {
			p.errChan <- err
			continue
		}

		oneRecord.InvertBit, err = readInt(row[14])
		if err != nil {
			p.errChan <- err
			continue
		}

		allRecords = append(allRecords, oneRecord)
	}

	return allRecords, nil
}

func (p *Parser) WriteDataToFile(ctx context.Context, records []database.Record) error {

	err := p.outFile.WriteData(ctx, records)
	if err != nil {
		return err
	}

	return nil
}

func readInt(str string) (int, error) {
	str = strings.TrimSpace(str)
	if str == "" {
		return 0, nil
	}

	val, err := strconv.Atoi(str)
	if err != nil {
		return 0, err
	}

	return val, nil
}

func readString(str string) (string, error) {
	str = strings.TrimSpace(str)

	return str, nil
}

func readBytes(str string) ([]byte, error) {
	str = strings.TrimSpace(str)

	return []byte(str), nil
}

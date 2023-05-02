package directory

import (
	"context"
	"io/ioutil"
	"time"

	"test_task/internal/app/config"
	"test_task/internal/app/database"
)

type FilesDirectory struct {
	path           string
	delay          time.Duration
	queue          chan string
	db             database.IDatabase
	processedFiles map[string]struct{}

	errChan chan error
}

func New(ctx context.Context, cfg *config.Config, queue chan string, db database.IDatabase, errChan chan error) (*FilesDirectory, error) {
	dir := FilesDirectory{}

	dir.path = cfg.FilesDirectory.FilesDirectory
	dir.delay = time.Millisecond * time.Duration(cfg.FilesDirectory.Delay)
	dir.queue = queue
	dir.db = db
	dir.errChan = errChan

	dbFiles, err := db.GetProcessedFiles(ctx)
	if err != nil {
		return nil, err
	}

	procFiles := make(map[string]struct{}, len(dbFiles))
	for _, file := range dbFiles {
		procFiles[file] = struct{}{}
	}
	dir.processedFiles = procFiles

	return &dir, nil
}

func (d *FilesDirectory) Run(ctx context.Context) {
	for {
		dirFiles, err := ioutil.ReadDir(d.path)
		if err != nil {
			d.errChan <- err
		}

		for _, file := range dirFiles {
			filePath := d.path + "\\" + file.Name()
			if _, ok := d.processedFiles[filePath]; !ok {
				d.queue <- filePath
				err = d.db.AddProcessedFile(ctx, filePath)
				if err != nil {
					d.errChan <- err
				}
				d.processedFiles[filePath] = struct{}{}
			}
		}

		time.Sleep(d.delay)
	}
}

package config

import (
	"os"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

type Config struct {
	App struct {
		QueueMaxSize int `env:"QUEUE_MAX_SIZE"`
	}
	Database struct {
		Host         string `env:"DB_HOST"`
		Port         int    `env:"DB_PORT"`
		User         string `env:"DB_USER"`
		Password     string `env:"DB_PASSWORD"`
		DatabaseName string `env:"DB_NAME"`
	}
	FilesDirectory struct {
		FilesDirectory string `env:"FILES_DIRECTORY"`
		Delay          int    `env:"CHECK_FILES_DIRECTORY_DELAY"`
	}
	Parser struct {
		OutFilesDirectory string  `env:"OUT_FILE_DIRECTORY"`
		PdfApiKey         string  `env:"PDF_API_KEY"`
		Font              string  `env:"STANDARD_FONT"`
		BoldFont          string  `env:"STANDARD_BOLD_FONT"`
		TableBorderWidth  float64 `env:"TABLE_BORDER_WIDTH"`
	}
	Service struct {
		Port     int32 `env:"SERVICE_PORT"`
		PageSize int32 `env:"SERVICE_PAGE_SIZE"`
	}
}

func New() (*Config, error) {
	err := loadEnv()
	if err != nil {
		return nil, err
	}

	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func loadEnv() error {
	err := godotenv.Load(os.Getenv("ENV_FILE"))
	if err != nil {
		return err
	}
	return nil
}

package config

import (
	"os"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

type Config struct {
	Database       DB
	FilesDirectory FilesDirectory
}

type DB struct {
	Host         string `env:"DB_HOST"`
	Port         int    `env:"DB_PORT"`
	User         string `env:"DB_USER"`
	Password     string `env:"DB_PASSWORD"`
	DatabaseName string `env:"DB_NAME"`
}

type FilesDirectory struct {
	FilesDirectory string `env:"FILES_DIRECTORY"`
	Delay          int    `env:"CHECK_FILES_DIRECTORY_DELAY"`
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

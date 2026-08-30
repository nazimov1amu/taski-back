package config

import (
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type AppConfig struct {
	Address string `env:"ADDRESS"`
	DatabaseDSN string `env:"DATABASE_DSN"`
}

func NewAppConfig() (*AppConfig, error) {
	var cfg AppConfig
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
package config

import (
	"github.com/caarlos0/env"
)

type Config struct {
	Listen   		string `env:"LISTEN" envDefault:"127.0.0.1:8080"`
	SearchingFolder string `env:"SDIR" envDefault:"./search"`
	LogLevel 		string `env:"LOG_LEVEL" envDefault:"info"`
}

func Load() (Config, error) {
	cfg := Config{}
	err := env.Parse(&cfg)
	return cfg, err
}

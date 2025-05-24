package config

import "github.com/caarlos0/env"

type Config struct {
	RedisURL string `env:"REDIS_URL,required"`
}

func New() (*Config, error) {
	var cfg Config

	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

package config

import "github.com/caarlos0/env"

type Config struct {
	RedisURL               string `env:"REDIS_URL,required"`
	SentimentIngestedTopic string `env:"SENTIMENT_INGESTED_TOPIC,required"`
	RestServerHost         string `env:"REST_SERVER_HOST,required"`
	RestServerPort         string `env:"REST_SERVER_PORT,required"`
}

func New() (*Config, error) {
	var cfg Config

	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

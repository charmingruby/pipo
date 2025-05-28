package config

import "github.com/caarlos0/env"

type Config struct {
	RestServerHost         string `env:"REST_SERVER_HOST,required"`
	RestServerPort         string `env:"REST_SERVER_PORT,required"`
	DatabaseHost           string `env:"DATABASE_HOST,required"`
	DatabasePort           string `env:"DATABASE_PORT,required"`
	DatabaseUser           string `env:"DATABASE_USER,required"`
	DatabasePassword       string `env:"DATABASE_PASSWORD,required"`
	DatabaseName           string `env:"DATABASE_NAME,required"`
	DatabaseSSL            string `env:"DATABASE_SSL,required"`
	RedisURL               string `env:"REDIS_URL,required"`
	SentimentIngestedTopic string `env:"SENTIMENT_INGESTED_TOPIC,required"`
}

func New() (*Config, error) {
	var cfg Config

	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

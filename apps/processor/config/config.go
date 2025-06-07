// Package config provides the configuration for the worker service.
package config

import "github.com/caarlos0/env"

// Config is the configuration for the worker service.
type Config struct {
	// RestServerPort is the port of the REST server.
	RestServerPort string `env:"REST_SERVER_PORT,required"`
	// DatabaseURL is the URL of the database.
	DatabaseURL string `env:"DATABASE_URL,required"`
	// RedisURL is the URL of the Redis.
	RedisURL string `env:"REDIS_URL,required"`
	// SentimentIngestedTopic is the topic of the sentiment ingested.
	SentimentIngestedTopic string `env:"SENTIMENT_INGESTED_TOPIC,required"`
}

// New constructs a new Config.
//
// Returns a new Config and an error if the configuration is invalid.
func New() (*Config, error) {
	var cfg Config

	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

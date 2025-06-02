// Package config provides the configuration for the worker service.
package config

import "github.com/caarlos0/env"

// Config is the configuration for the worker service.
type Config struct {
	// RestServerPort is the port of the REST server.
	RestServerPort string `env:"REST_SERVER_PORT,required"`
	// DatabaseHost is the host of the database.
	DatabaseHost string `env:"DATABASE_HOST,required"`
	// DatabasePort is the port of the database.
	DatabasePort string `env:"DATABASE_PORT,required"`
	// DatabaseUser is the user of the database.
	DatabaseUser string `env:"DATABASE_USER,required"`
	// DatabasePassword is the password of the database.
	DatabasePassword string `env:"DATABASE_PASSWORD,required"`
	// DatabaseName is the name of the database.
	DatabaseName string `env:"DATABASE_NAME,required"`
	// DatabaseSSL is the SSL of the database.
	DatabaseSSL string `env:"DATABASE_SSL,required"`
	// RedisURL is the URL of the Redis server.
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

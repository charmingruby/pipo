package main

import (
	"os"

	"github.com/charmingruby/pipo/lib/broker/redis"
	"github.com/charmingruby/pipo/lib/logger"
	"github.com/charmingruby/pipo/lib/persistence/postgres"
	"github.com/charmingruby/pipo/service/worker/config"
	"github.com/charmingruby/pipo/service/worker/internal/core/service"
	"github.com/charmingruby/pipo/service/worker/internal/database/repository"
	"github.com/charmingruby/pipo/service/worker/internal/delivery/event"
	"github.com/joho/godotenv"
)

func main() {
	logger := logger.New()

	if err := godotenv.Load(); err != nil {
		logger.Warn("failed to find .env file", "error", err)
	}

	cfg, err := config.New()
	if err != nil {
		logger.Error("failed to load config", "error", err)

		os.Exit(1)
	}

	logger.Info("config loaded")

	redisClient, err := redis.NewClient(cfg.RedisURL)
	if err != nil {
		logger.Error("failed to connect to redis", "error", err)

		os.Exit(1)
	}

	logger.Info("redis connected")

	redisBroker := redis.NewStream(redisClient)

	db, err := postgres.New(logger, postgres.ConnectionInput{
		Host:         cfg.DatabaseHost,
		Port:         cfg.DatabasePort,
		User:         cfg.DatabaseUser,
		Password:     cfg.DatabasePassword,
		DatabaseName: cfg.DatabaseName,
		SSL:          cfg.DatabaseSSL,
	})
	if err != nil {
		logger.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}

	logger.Info("database connected")

	repo, err := repository.NewPostgresSentimentRepository(db.Conn)
	if err != nil {
		logger.Error("failed to create postgres sentiment repository", "error", err)

		os.Exit(1)
	}

	service := service.New(logger, redisBroker, repo, cfg.SentimentIngestedTopic)

	eventHandler := event.New(logger, redisBroker, event.TopicInput{
		SentimentIngested: cfg.SentimentIngestedTopic,
	}, service)

	if err := eventHandler.Subscribe(); err != nil {
		logger.Error("failed to subscribe to event handler", "error", err)

		os.Exit(1)
	}
}

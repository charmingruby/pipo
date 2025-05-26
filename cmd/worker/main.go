package main

import (
	"os"

	"github.com/charmingruby/pipo/config"
	"github.com/charmingruby/pipo/internal/sentiment/core/service"
	"github.com/charmingruby/pipo/internal/sentiment/delivery/event"
	"github.com/charmingruby/pipo/internal/shared/broker"
	"github.com/charmingruby/pipo/pkg/logger"
	"github.com/charmingruby/pipo/pkg/redis"
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

	redisClient, err := redis.New(cfg.RedisURL)
	if err != nil {
		logger.Error("failed to connect to redis", "error", err)

		os.Exit(1)
	}

	logger.Info("redis connected")

	redisBroker := broker.NewRedisStream(redisClient)

	service := service.New(logger, redisBroker, cfg.SentimentIngestedTopic)

	eventHandler := event.New(redisBroker, event.TopicInput{
		SentimentIngested: cfg.SentimentIngestedTopic,
	}, service)

	if err := eventHandler.Subscribe(); err != nil {
		logger.Error("failed to subscribe to event handler", "error", err)

		os.Exit(1)
	}
}

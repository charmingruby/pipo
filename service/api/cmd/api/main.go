package main

import (
	"os"

	"github.com/charmingruby/pipo/lib/broker/redis"
	"github.com/charmingruby/pipo/lib/http/rest"
	"github.com/charmingruby/pipo/lib/logger"
	"github.com/charmingruby/pipo/service/api/config"
	"github.com/charmingruby/pipo/service/api/internal/core/service"
	"github.com/charmingruby/pipo/service/api/internal/delivery/http/rest/endpoint"
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

	service := service.New(logger, redisBroker, cfg.SentimentIngestedTopic)

	server, router := rest.New(cfg.RestServerHost, cfg.RestServerPort)

	endpoint := endpoint.New(router, service)
	endpoint.Register()

	if err := server.Start(); err != nil {
		logger.Error("failed to start rest server", "error", err)

		os.Exit(1)
	}
}

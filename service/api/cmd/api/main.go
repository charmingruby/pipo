package main

import (
	"context"
	"os"
	"os/signal"
	"time"

	"github.com/charmingruby/pipo-lib/broker/redis"
	"github.com/charmingruby/pipo-lib/http/rest"
	"github.com/charmingruby/pipo-lib/logger"
	"github.com/charmingruby/pipo/service/api/config"
	"github.com/charmingruby/pipo/service/api/internal/core/service"
	"github.com/charmingruby/pipo/service/api/internal/delivery/http/rest/endpoint"
	"github.com/charmingruby/pipo/service/api/internal/platform/monitoring/health"
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

	logger.Info("redis broker created")

	service := service.New(logger, redisBroker, cfg.SentimentIngestedTopic)

	server, router := rest.New(cfg.RestServerHost, cfg.RestServerPort)

	endpoint := endpoint.New(router, service)
	endpoint.Register()

	health.
		NewHealth(router, logger, redisClient).
		RegisterProbes()

	logger.Info("registered probes")

	logger.Info("rest server started", "port", cfg.RestServerPort)

	go func() {
		if err := server.Start(); err != nil {
			logger.Error("failed to start rest server", "error", err)

			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	logger.Info("received interrupt signal")

	gracefulShutdown(logger, server, redisClient)
}

func gracefulShutdown(
	logger *logger.Logger,
	server *rest.Server,
	redis *redis.Client,
) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Stop(ctx); err != nil {
		logger.Error("failed to stop rest server", "error", err)
	}

	logger.Info("rest server stopped")

	if err := redis.Close(); err != nil {
		logger.Error("failed to close redis client", "error", err)
	}

	logger.Info("redis client closed")
}

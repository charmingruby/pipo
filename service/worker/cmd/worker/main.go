package main

import (
	"context"
	"os"
	"os/signal"
	"time"

	"github.com/charmingruby/pipo-lib/broker/redis"
	"github.com/charmingruby/pipo-lib/http/rest"
	"github.com/charmingruby/pipo-lib/logger"
	"github.com/charmingruby/pipo-lib/persistence/postgres"
	"github.com/charmingruby/pipo/service/worker/config"
	"github.com/charmingruby/pipo/service/worker/internal/core/service"
	"github.com/charmingruby/pipo/service/worker/internal/database/repository"
	"github.com/charmingruby/pipo/service/worker/internal/delivery/event"
	"github.com/charmingruby/pipo/service/worker/internal/platform/monitoring/health"
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

	repo := repository.NewPostgresSentimentRepository(db.Conn)

	service := service.New(logger, redisBroker, repo, cfg.SentimentIngestedTopic)

	go func() {
		eventHandler := event.New(logger, redisBroker, event.TopicInput{
			SentimentIngested: cfg.SentimentIngestedTopic,
		}, service)

		logger.Info("event handler subscribed")

		if err := eventHandler.Subscribe(); err != nil {
			logger.Error("failed to subscribe to event handler", "error", err)

			os.Exit(1)
		}
	}()

	server, router := rest.New(cfg.RestServerHost, cfg.RestServerPort)

	health.
		NewHealth(router, logger, db, redisClient).
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

	gracefulShutdown(logger, redisClient, db, server)
}

func gracefulShutdown(
	logger *logger.Logger,
	redis *redis.Client,
	db *postgres.Client,
	server *rest.Server,
) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Stop(ctx); err != nil {
		logger.Error("failed to stop rest server", "error", err)
	}

	if err := db.Close(ctx); err != nil {
		logger.Error("failed to close database", "error", err)
	}

	logger.Info("database closed")

	if err := redis.Close(); err != nil {
		logger.Error("failed to close redis client", "error", err)
	}

	logger.Info("redis client closed")
}

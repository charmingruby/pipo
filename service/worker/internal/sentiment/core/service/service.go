package service

import (
	"github.com/charmingruby/pipo/internal/sentiment/core/repository"
	"github.com/charmingruby/pipo/internal/shared/broker"
	"github.com/charmingruby/pipo/lib/logger"
)

type Service struct {
	logger               *logger.Logger
	broker               broker.Broker
	sentimentRepository  repository.SentimentRepository
	sentimentIngestTopic string
}

func New(logger *logger.Logger, broker broker.Broker, sentimentRepo repository.SentimentRepository, sentimentIngestTopic string) *Service {
	return &Service{
		logger:               logger,
		broker:               broker,
		sentimentRepository:  sentimentRepo,
		sentimentIngestTopic: sentimentIngestTopic,
	}
}

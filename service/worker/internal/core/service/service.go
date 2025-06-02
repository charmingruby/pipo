// Package service provides the business logic for the application.
package service

import (
	"github.com/charmingruby/pipo-lib/broker"
	"github.com/charmingruby/pipo-lib/logger"
	"github.com/charmingruby/pipo/service/worker/internal/core/repository"
)

// Service is the service for the application.
type Service struct {
	logger               *logger.Logger
	broker               broker.Broker
	sentimentRepository  repository.SentimentRepository
	sentimentIngestTopic string
}

// New constructs a new Service.
//
// logger is the logger for the service.
// broker is the broker for the service.
// sentimentRepo is the sentiment repository for the service.
// sentimentIngestTopic is the topic of the sentiment ingested.
//
// Returns a new Service.
func New(
	logger *logger.Logger,
	broker broker.Broker,
	sentimentRepo repository.SentimentRepository,
	sentimentIngestTopic string,
) *Service {
	return &Service{
		logger:               logger,
		broker:               broker,
		sentimentRepository:  sentimentRepo,
		sentimentIngestTopic: sentimentIngestTopic,
	}
}

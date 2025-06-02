// Package service provides the business logic for the application.
package service

import (
	"github.com/charmingruby/pipo-lib/broker"
	"github.com/charmingruby/pipo-lib/logger"
)

// Service is the business logic for the application.
type Service struct {
	logger               *logger.Logger
	broker               broker.Broker
	sentimentIngestTopic string
}

// New constructs a new Service.
//
// logger is the logger for the service.
// broker is the broker for the service.
// sentimentIngestTopic is the topic of the sentiment ingested.
//
// Returns a new Service.
func New(logger *logger.Logger, broker broker.Broker, sentimentIngestTopic string) *Service {
	return &Service{
		logger:               logger,
		broker:               broker,
		sentimentIngestTopic: sentimentIngestTopic,
	}
}

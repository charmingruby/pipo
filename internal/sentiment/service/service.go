package service

import (
	"github.com/charmingruby/pipo/internal/shared/messaging"
	"github.com/charmingruby/pipo/pkg/logger"
)

type Service struct {
	logger               *logger.Logger
	broker               messaging.Broker
	sentimentIngestTopic string
}

func NewService(logger *logger.Logger, broker messaging.Broker, sentimentIngestTopic string) *Service {
	return &Service{
		logger:               logger,
		broker:               broker,
		sentimentIngestTopic: sentimentIngestTopic,
	}
}

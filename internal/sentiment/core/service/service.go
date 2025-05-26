package service

import (
	"github.com/charmingruby/pipo/internal/shared/broker"
	"github.com/charmingruby/pipo/pkg/logger"
)

type Service struct {
	logger               *logger.Logger
	broker               broker.Broker
	sentimentIngestTopic string
}

func New(logger *logger.Logger, broker broker.Broker, sentimentIngestTopic string) *Service {
	return &Service{
		logger:               logger,
		broker:               broker,
		sentimentIngestTopic: sentimentIngestTopic,
	}
}
